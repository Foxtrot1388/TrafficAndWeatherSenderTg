package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/sender/telegram"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic"
	"github.com/go-co-op/gocron"
)

type trafficGeter interface {
	InfoYandex(url string, w io.Writer)
}

type sender interface {
	SendPhoto(name, text string, bytes []byte) error
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := config.Get()
	trafficya := traffic.New()
	sendertg, err := telegram.New(cfg.Token, cfg.ChatID)
	if err != nil {
		log.Fatal(err)
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At(cfg.Time).Do(func() { do(cfg, trafficya, sendertg) })

}

func do(cfg *config.Config, traffic trafficGeter, sendermes sender) {

	done := make(chan struct{})

	go func() {
		sendTrafficInfo(cfg, traffic, sendermes)
		done <- struct{}{}
	}()

	<-done

}

func sendTrafficInfo(cfg *config.Config, traffic trafficGeter, sendermes sender) {

	buf := bytes.NewBuffer([]byte{})
	traffic.InfoYandex(cfg.URL, buf)
	if buf.Len() == 0 {
		log.Println("Failed to get traffic information!")
		return
	}

	curtime := time.Now()
	fileName := fmt.Sprintf("traffic %s.png", curtime.Format("2006-01-02"))
	text := fmt.Sprintf("Пробки %s", curtime.Format("2006-01-02 3:4:5"))
	if err := sendermes.SendPhoto(fileName, text, buf.Bytes()); err != nil {
		log.Fatal(err)
	}

}
