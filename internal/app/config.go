package app

import (
	"github.com/caarlos0/env/v6"
	"log"
)

var ServerURL = "localhost:8080"
var BaseURL = "http://localhost:8080"

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func configServerEnv() {
	var cfg = Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	ServerURL = cfg.ServerAddr
	BaseURL = cfg.BaseURL
}

//
