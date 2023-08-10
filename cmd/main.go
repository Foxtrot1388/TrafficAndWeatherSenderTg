package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/sender/telegram"
	task "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/task/Google"
	traffic "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic/Yandex"
	weather "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/weather/OpenMeteo"
	"github.com/go-co-op/gocron"
)

type trafficGeter interface {
	Info(url string, w io.Writer)
}

type weatherGeter interface {
	Info(ctx context.Context, timezone string, lat float64, lon float64, filter func(time.Time) bool) (*weather.WeatherInfo, error)
}

type taskGeter interface {
	Info() ([]task.TaskResult, error)
}

type sender interface {
	SendPhoto(name, text string, bytes []byte) error
	SendText(text string) error
}

var errorFailWeather = errors.New("failed to get weather information")
var errorFailTraffic = errors.New("failed to get traffic information")
var errorFailTask = errors.New("failed to get task information")

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := config.Get()
	trafficya := traffic.New()
	weatherom, err := weather.New()
	if err != nil {
		log.Fatal(err)
	}
	taskg, err := task.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	sendertg, err := telegram.New(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatal(err)
	}

	// TODO context

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At(cfg.TimeToSend).Do(func() { do(cfg, trafficya, weatherom, taskg, sendertg) })

}

func do(cfg *config.Config, traffic trafficGeter, weather weatherGeter, task taskGeter, sendermes sender) {

	log.Println("The time has come")
	done := make(chan struct{})

	go func() {
		err := sendTrafficInfo(cfg, traffic, sendermes)
		if err != nil {
			if errors.Is(err, errorFailTraffic) {
				log.Println(err.Error())
			} else {
				log.Fatal(err)
			}
		}
		done <- struct{}{}
	}()

	go func() {
		err := sendWeatherInfo(cfg, weather, sendermes)
		if err != nil {
			if errors.Is(err, errorFailWeather) {
				log.Println(err.Error())
			} else {
				log.Fatal(err)
			}
		}
		done <- struct{}{}
	}()

	go func() {
		err := sendTaskInfo(cfg, task, sendermes)
		if err != nil {
			if errors.Is(err, errorFailTask) {
				log.Println(err.Error())
			} else {
				log.Fatal(err)
			}
		}
		done <- struct{}{}
	}()

	<-done
	<-done
	<-done

}

func sendTrafficInfo(cfg *config.Config, traffic trafficGeter, sendermes sender) error {

	buf := bytes.NewBuffer([]byte{})
	traffic.Info(cfg.Traffic.URL, buf)
	if buf.Len() == 0 {
		return errorFailTraffic
	}

	curtime := time.Now()
	fileName := fmt.Sprintf("traffic %s.png", curtime.Format("2006-01-02"))
	text := fmt.Sprintf("Пробки %s", curtime.Format("2006-01-02 3:4:5"))
	if err := sendermes.SendPhoto(fileName, text, buf.Bytes()); err != nil {
		return err
	}

	return nil

}

func sendWeatherInfo(cfg *config.Config, weather weatherGeter, sendermes sender) error {

	filter, err := filterWeather(cfg)
	if err != nil {
		return err
	}

	result, err := weather.Info(context.Background(), cfg.Weather.Timezone, cfg.Weather.Lat, cfg.Weather.Lon, filter)
	if err != nil {
		return errorFailWeather
	}

	if err = sendermes.SendText(result.String()); err != nil {
		return err
	}

	return nil

}

func filterWeather(cfg *config.Config) (func(time.Time) bool, error) {

	curday := time.Now().Day()
	timesplit := strings.Split(cfg.Weather.Time, "-")
	timeAt, err := strconv.Atoi(strings.Split(timesplit[0], ":")[0])
	if err != nil {
		return nil, err
	}
	timeTo, err := strconv.Atoi(strings.Split(timesplit[1], ":")[0])
	if err != nil {
		return nil, err
	}

	filter := func(t time.Time) bool {
		return t.Day() == curday && t.Hour() >= timeAt && t.Hour() <= timeTo
	}

	return filter, nil

}

func sendTaskInfo(cfg *config.Config, task taskGeter, sendermes sender) error {

	result, err := task.Info()
	if err != nil {
		return errorFailTask
	}

	if len(result) > 0 {
		sb := strings.Builder{}
		sb.WriteString("Список задач на сегодня:")
		for i := 0; i < len(result); i++ {
			sb.WriteString("\r\n")
			sb.WriteString(result[i].String())
		}
		if err = sendermes.SendText(sb.String()); err != nil {
			return err
		}
	}

	return nil

}
