package app

import (
	"fmt"
	"os"
	"strconv"
)

var ServerURL = "localhost:8080"
var BaseURL = "http://localhost:8080"

type Config struct {
	port int
	host string
}

func (c *Config) ConfigServerEnv() {

	//flag.IntVar(&c.port, "port", 8080, "port to listen on")
	//flag.StringVar(&c.host, "host", "localhost", "host to listen on")
	//flag.Parse()
	var err error
	c.port, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil || c.port < 1025 || c.port > 65535 {
		c.port = 8080
	}
	c.host = os.Getenv("SERVER_HOST")
	if c.host == "" {
		c.host = "localhost"
	}

	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	fmt.Printf("server start on %s\n", addr)

	ServerURL = addr
	BaseURL = "http://" + addr
}
