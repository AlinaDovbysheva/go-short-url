package app

import (
	"flag"
	"fmt"
)

var ServerURL = ":8080"
var BaseURL = "http://localhost:8080"
var FilePath = "db1.log"

type Config struct {
	port     string
	host     string
	filepath string
}

func (c *Config) ConfigServerEnv() {

	flag.StringVar(&c.port, "a", ":8080", "port to listen on")
	flag.StringVar(&c.host, "b", "http://localhost:8080", "host to listen on")
	flag.StringVar(&c.filepath, "f", "db1.log", "file path to storage")
	flag.Parse()

	ServerURL = c.port
	BaseURL = c.host
	FilePath = c.filepath

	fmt.Printf("server start on %s ", ServerURL)

	/*
		ServerURL = os.Getenv("SERVER_ADDRESS")
		if ServerURL == "" {
			ServerURL = ":8080"
		}
		BaseURL = os.Getenv("BASE_URL")
		if BaseURL == "" {
			BaseURL = "http://localhost:8080"
		}
		FilePath = os.Getenv("FILE_STORAGE_PATH")
		if FilePath == "" {
			FilePath = "URLdb.log"
		}
	*/
}
