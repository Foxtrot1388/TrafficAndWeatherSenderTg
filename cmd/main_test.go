package main

import (
	"testing"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/sender/telegram"
	traffic "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic/Yandex"
	weather "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/weather/OpenMeteo"
)

func TestTrafficInfo(t *testing.T) {

	cfg := config.Get()
	trafficya := traffic.New()
	sendertg, err := telegram.New(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		t.Error(err)
	}

	err = sendTrafficInfo(cfg, trafficya, sendertg)
	if err != nil {
		t.Error(err)
	}

}

func TestWeatherInfo(t *testing.T) {

	cfg := config.Get()
	weatherow, err := weather.New()
	if err != nil {
		t.Error(err)
	}

	sendertg, err := telegram.New(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		t.Error(err)
	}

	err = sendWeatherInfo(cfg, weatherow, sendertg)
	if err != nil {
		t.Error(err)
	}

}
