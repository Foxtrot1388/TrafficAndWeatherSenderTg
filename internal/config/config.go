package config

import (
	"flag"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TimeToSend string `yaml:"time_to_send"`
	Telegram   struct {
		ChatID int64  `yaml:"chat_id"`
		Token  string `yaml:"token_id"`
	} `yaml:"telegram"`
	Traffic struct {
		Enable bool   `yaml:"enable"`
		URL    string `yaml:"url"`
	} `yaml:"traffic"`
	Weather struct {
		Enable   bool    `yaml:"enable"`
		Lat      float64 `yaml:"lat"`
		Lon      float64 `yaml:"lon"`
		Time     string  `yaml:"time"`
		Timezone string  `yaml:"timezone"` // https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	} `yaml:"weather"`
	Task struct {
		Enable       bool   `yaml:"enable"`
		Caledndarid  string `yaml:"calendarid"`
		Clientsecret string `yaml:"clientsecret"`
		Token        string `yaml:"token"`
	} `yaml:"task"`
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
