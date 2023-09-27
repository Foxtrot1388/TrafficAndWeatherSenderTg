package main

import (
	"context"
	"testing"

	"github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/config"
	task "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/task/Google"
	traffic "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/traffic/Yandex"
	weather "github.com/Foxtrot1388/TrafficAndWeatherSenderTg/internal/weather/OpenMeteo"
)

type sendermok struct {
	t *testing.T
}

func (s *sendermok) SendPhoto(name, text string, bytes []byte) error {
	s.t.Log(name, text)
	return nil
}

func (s *sendermok) SendText(text string) error {
	s.t.Log(text)
	return nil
}

func TestTrafficInfo(t *testing.T) {

	cfg := config.Get()

	trafficya := traffic.New()
	sendertg := &sendermok{t: t}

	err := sendTrafficInfo(context.Background(), cfg, trafficya, sendertg)
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
	sendertg := &sendermok{t: t}

	err = sendWeatherInfo(context.Background(), cfg, weatherow, sendertg)
	if err != nil {
		t.Error(err)
	}

}

func TestCalendar(t *testing.T) {

	cfg := config.Get()

	taskg, err := task.New(context.Background(), cfg.Task.Token, cfg.Task.Clientsecret)
	if err != nil {
		t.Error(err)
	}
	sendertg := &sendermok{t: t}

	err = sendTaskInfo(context.Background(), cfg, taskg, sendertg)
	if err != nil {
		t.Error(err)
	}

}
