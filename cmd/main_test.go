package main

import (
	"testing"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/sender/telegram"
	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic"
)

func TestDoAll(t *testing.T) {

	cfg := config.Get()
	trafficya := traffic.New()
	sendertg, err := telegram.New(cfg.Token, cfg.ChatID)
	if err != nil {
		t.Error(err)
	}

	do(cfg, trafficya, sendertg)

}
