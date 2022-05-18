package app

import (
	"flag"
	"fmt"
	"os"
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

	flag.StringVar(&c.port, "a", ServerURL, "port to listen on")
	flag.StringVar(&c.host, "b", BaseURL, "host to listen on")
	flag.StringVar(&c.filepath, "f", FilePath, "file path to storage")
	flag.Parse()

	ServerURL = c.port
	BaseURL = c.host //+ c.port
	FilePath = c.filepath

	fmt.Printf("server start on %s ", ServerURL)

}
