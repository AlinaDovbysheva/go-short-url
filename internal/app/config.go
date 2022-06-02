package app

import (
	"database/sql"
	"flag"
	"os"
)

var ServerURL = ":8080"
var BaseURL = "http://localhost:8080"
var FilePath = "db1.log"
var DatabaseDsn = "host=localhost port=5432 user=postgres password=plokij098 dbname=DB_shortner sslmode=disable"

//"postgres://user:pass@localhost/bookstore"
type Config struct {
	port        string
	host        string
	filepath    string
	databasedsn string
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
	/*if FilePath == "" {
		FilePath = "" //"URLdb.log"
	}*/

	DatabaseDsn = os.Getenv("DATABASE_DSN")
	if DatabaseDsn == "" {
		DatabaseDsn = "host=localhost port=5432 user=postgres password=plokij098 dbname=DB_shortner sslmode=disable"
	}

	flag.StringVar(&c.port, "a", ServerURL, "port to listen on")
	flag.StringVar(&c.host, "b", BaseURL, "host to listen on")
	flag.StringVar(&c.filepath, "f", FilePath, "file path to storage")
	flag.StringVar(&c.databasedsn, "d", DatabaseDsn, "database to connect")
	flag.Parse()

	ServerURL = c.port
	BaseURL = c.host //+ c.port
	FilePath = c.filepath
	DatabaseDsn = c.databasedsn
	//fmt.Printf("server start on %s  file set to %s /n", ServerURL, FilePath)
}

func (c *Config) PingDB() error {
	db, err := sql.Open("postgres", c.databasedsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}
