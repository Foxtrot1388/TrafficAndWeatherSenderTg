package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
)

func main() {

	cfg := config.Get()
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At(cfg.Time).Do(do)
}

func do() {

	cfg := config.Get()
	curtime := time.Now()

	buf := bytes.NewBuffer([]byte{})
	traffic.GetInfo(cfg.URL, buf)

	bot, err := NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	b := tgbotapi.FileBytes{Name: fmt.Sprintf("traffic %s.png", curtime.Format("2006-01-02")), Bytes: buf.Bytes()}

	msg := tgbotapi.NewPhoto(cfg.ChatID, b)
	msg.Caption = fmt.Sprintf("Пробки %s", curtime.Format("2006-01-02 3:4:5"))
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

}

func NewBotAPI(token string) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPIWithClient(token, APIEndpoint, &http.Client{})
}
