package storage

import (
	"errors"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

type DBurl interface {
	GetUrl(shortUrl string) (string, error)
	PutUrl(inputUrl string) (string, error)
}

func GetUrl(shortUrl string) (string, error) {
	sId := mapUrl[shortUrl]
	if sId == "" {
		return "", errors.New("Id is absent in DB")
	}
	return sId, nil
}

func PutUrl(inputUrl string) (string, error) {
	return WriteURL(inputUrl), nil
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
