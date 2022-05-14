package app

import (
	"os"
)

var ServerURL = ":8080"
var BaseURL = "http://localhost:8080"

type Config struct {
	port int
	host string
}

func (c *Config) ConfigServerEnv() {

	/*flag.IntVar(&c.port, "port", 8080, "port to listen on")
	flag.StringVar(&c.host, "host", "localhost", "host to listen on")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	fmt.Printf("server start on %s\n", addr)

	ServerURL = addr
	BaseURL = "http://" + addr*/

	ServerURL = os.Getenv("SERVER_ADDRESS")
	BaseURL = os.Getenv("BASE_URL")

}
