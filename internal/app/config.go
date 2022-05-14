package app

import (
	"flag"
	"fmt"
)

var ServerURL = "localhost:8080"
var BaseURL = "http://localhost:8080"

type Config struct {
	port int
	host string
}

func (c *Config) ConfigServerEnv() {

	flag.IntVar(&c.port, "SERVER_PORT", 8080, "port to listen on")
	flag.StringVar(&c.host, "SERVER_HOST", "localhost", "host to listen on")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	fmt.Printf("server start on %s\n", addr)

	ServerURL = addr
	BaseURL = "http://" + addr
}
