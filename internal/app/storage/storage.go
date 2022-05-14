package storage

import (
	"errors"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

type DBurl interface {
	GetURL(shortURL string) (string, error)
	PutURL(inputURL string) (string, error)
}

func GetURL(shortURL string) (string, error) {
	sID := mapURL[shortURL]
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func PutURL(inputURL string) (string, error) {
	return WriteURL(inputURL), nil
}

var mapURL = make(map[string]string)

func WriteURL(url string) (id string) {
	id = FindURLKey(url)
	if id == "" {
		id = util.RandStringBytes(7)
		mapURL[id] = url
	}
	return id
}

func FindURL(key string) (result string) {
	return mapURL[key]
}

func FindURLKey(url string) (key string) {
	key = ""
	for k, v := range mapURL {
		if v == url {
			return k
		}
	}
	return key
}
