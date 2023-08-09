package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/sender/telegram"
	traffic "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic/Yandex"
	weather "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/weather/OpenMeteo"
	"github.com/go-co-op/gocron"
)

type trafficGeter interface {
	Info(url string, w io.Writer)
}

type weatherGeter interface {
	Info(ctx context.Context, timezone string, lat, lon float64) error
}

type sender interface {
	SendPhoto(name, text string, bytes []byte) error
	SendText(text string) error
}

var errorFailWeather = errors.New("failed to get weather information")
var errorFailTraffic = errors.New("failed to get traffic information")

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := config.Get()
	trafficya := traffic.New()
	weatherom, err := weather.New()
	if err != nil {
		log.Fatal(err)
	}
	sendertg, err := telegram.New(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatal(err)
	}

	// TODO context

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At(cfg.TimeToSend).Do(func() { do(cfg, trafficya, weatherom, sendertg) })

}

func do(cfg *config.Config, traffic trafficGeter, weather weatherGeter, sendermes sender) {

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

	if err := weather.Info(context.Background(), cfg.Weather.Timezone, cfg.Weather.Lat, cfg.Weather.Lon); err != nil {
		return errorFailWeather
	}

	return nil

}
