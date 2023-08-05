package config

import (
	"flag"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	URL    string `yaml:"url"`
	Time   string `yaml:"time"`
	ChatID int64  `yaml:"chat_id"`
	Token  string `yaml:"token_id"`
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("../"+mustConfig(), instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}

func mustConfig() string {

	filename := flag.String(
		"cfg",
		"app.yaml",
		"config file",
	)

	flag.Parse()

	return *filename
}
