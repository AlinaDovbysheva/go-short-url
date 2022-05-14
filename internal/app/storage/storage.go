package storage

import (
	"errors"
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
	id = FindURL(url)
	if id == "" {
		id = util.RandStringBytes(7)
		mapUrl[id] = url
	}
	return id
}

func FindURL(url string) (result string) {
	result = ""
	for k, v := range mapUrl {
		if v == url {
			return k
		}
	}
	return result
}
