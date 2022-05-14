package storage

import (
	"errors"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

type DBurl interface {
	GetUrl(shortURL string) (string, error)
	PutUrl(inputURL string) (string, error)
}

func GetUrl(shortURL string) (string, error) {
	sID := mapUrl[shortURL]
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func PutURL(inputURL string) (string, error) {
	return WriteURL(inputURL), nil
}

var mapUrl = make(map[string]string)

func WriteURL(url string) (id string) {
	id = FindURLKey(url)
	if id == "" {
		id = util.RandStringBytes(7)
		mapUrl[id] = url
	}
	return id
}

func FindURL(key string) (result string) {
	return mapUrl[key]
}

func FindURLKey(url string) (key string) {
	key = ""
	fmt.Println(" mapUrl:")
	fmt.Println(mapUrl)
	for k, v := range mapUrl {
		if v == url {
			return k
		}
	}
	return key
}
