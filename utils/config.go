package utils

import (
	"time"
)

type TBot struct {
	Token            string        `json:"Token" yaml:"Token"`
	Channel          int64         `json:"Channel" yaml:"Channel"`
	Timeout          time.Duration `json:"-" yaml:"-"`
	SuccessImagePath string        `json:"SuccessImagePath" yaml:"SuccessImagePath"`
}

type Cache struct {
	Dir string `json:"Dir" yaml:"Dir"`
}

type App struct {
	Cache        Cache  `json:"Cache" yaml:"Cache"`
	OpenApiToken string `json:"OpenApiToken" yaml:"OpenApiToken"`
}

type Config struct {
	TBot TBot `json:"TBot" yaml:"TBot"`
	App  App  `json:"App" yaml:"App"`
}

func GetConfig() (*Config, error) {
	var c Config

	err := ReadFileAndUnmarshal("./config/config.yaml", &c)
	if err != nil {
		return nil, err
	}

	c.TBot.Timeout = time.Second * 10

	return &c, nil
}
