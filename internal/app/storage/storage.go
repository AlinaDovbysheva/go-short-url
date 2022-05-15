package storage

import (
	"errors"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

func NewInMap() DBurl {
	return &InMap{map[string]string{}}
}

type InMap struct {
	mapURL map[string]string
}

func (m *InMap) GetURL(shortURL string) (string, error) {
	sID := m.mapURL[shortURL]
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func (m *InMap) PutURL(inputURL string) (string, error) {
	id := ""
	for k, v := range m.mapURL {
		if v == inputURL {
			id = k
		}
	}
	if id == "" {
		id = util.RandStringBytes(7)
		m.mapURL[id] = inputURL
	}
	return id, nil
}

/*
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
*/
