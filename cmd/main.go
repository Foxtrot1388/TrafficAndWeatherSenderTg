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
	Info(ctx context.Context, url string, w io.Writer)
}

type weatherGeter interface {
	Info(ctx context.Context, loc *weather.Location, filter func(time.Time) bool) (*weather.WeatherInfo, error)
	NewLocation(lat, lon float64, timezone string) (*weather.Location, error)
}

type taskGeter interface {
	Info(ctx context.Context, calendarid string) ([]task.TaskResult, error)
}

type sender interface {
	SendPhoto(name, text string, bytes []byte) error
	SendText(text string) error
}

type application struct {
	Taskapp    taskGeter
	Senderapp  sender
	Weatherapp weatherGeter
	Trafficapp trafficGeter
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := config.Get()
	app := getApp(cfg)

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At(cfg.TimeToSend).Do(func() { do(context.Background(), cfg, app) })
	s.StartBlocking()

}

func getApp(cfg *config.Config) *application {

	trafficya := traffic.New()

	weatherom, err := weather.New()
	if err != nil {
		log.Fatal(err)
	}

	taskg, err := task.New(context.TODO(), cfg.Task.Token, cfg.Task.Clientsecret)
	if err != nil {
		log.Fatal(err)
	}

	sendertg, err := telegram.New(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatal(err)
	}

	return &application{Taskapp: taskg, Trafficapp: trafficya, Weatherapp: weatherom, Senderapp: sendertg}

}

func do(ctx context.Context, cfg *config.Config, app *application) {

	log.Println("The time has come")
	done := make(chan struct{})
	defer close(done)

	go func() {
		startDo(func() error { return sendTrafficInfo(ctx, cfg, app.Trafficapp, app.Senderapp) })
		done <- struct{}{}
	}()

	go func() {
		startDo(func() error { return sendWeatherInfo(ctx, cfg, app.Weatherapp, app.Senderapp) })
		done <- struct{}{}
	}()

	go func() {
		startDo(func() error { return sendTaskInfo(ctx, cfg, app.Taskapp, app.Senderapp) })
		done <- struct{}{}
	}()

	<-done
	<-done
	<-done

}

func startDo(job func() error) {
	err := job()
	if err != nil {
		log.Println(err.Error())
	}
}

func sendTrafficInfo(ctx context.Context, cfg *config.Config, traffic trafficGeter, sendermes sender) error {

	buf := bytes.NewBuffer([]byte{})
	traffic.Info(ctx, cfg.Traffic.URL, buf)
	if buf.Len() == 0 {
		return errors.New("failed to get traffic information")
	}

	curtime := time.Now()
	fileName := fmt.Sprintf("traffic %s.png", curtime.Format("2006-01-02"))
	text := fmt.Sprintf("Пробки %s", curtime.Format("2006-01-02 3:4:5"))
	if err := sendermes.SendPhoto(fileName, text, buf.Bytes()); err != nil {
		return err
	}

	return nil

}

func sendWeatherInfo(ctx context.Context, cfg *config.Config, weather weatherGeter, sendermes sender) error {

	filter, err := filterWeather(cfg)
	if err != nil {
		return err
	}

	loc, err := weather.NewLocation(cfg.Weather.Lat, cfg.Weather.Lon, cfg.Weather.Timezone)
	if err != nil {
		return err
	}

	result, err := weather.Info(ctx, loc, filter)
	if err != nil {
		return err
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

func sendTaskInfo(ctx context.Context, cfg *config.Config, task taskGeter, sendermes sender) error {

	result, err := task.Info(ctx, cfg.Task.Caledndarid)
	if err != nil {
		return err
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
