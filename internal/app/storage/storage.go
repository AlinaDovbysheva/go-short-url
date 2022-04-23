package storage

import "github.com/AlinaDovbysheva/go-short-url/internal/app/util"

type urlItem struct {
	ID  string
	URL string
}

var urls []urlItem

func WriteURL(url string) (id string) {
	id = FindURL(url)
	if id == "" {
		id = util.RandStringBytes(7)
		item := urlItem{ID: id, URL: url}
		urls = append(urls, item)
	}
	return id
}

func FindURL(url string) (result string) {
	result = ""
	for _, row := range urls {
		if row.URL == url {
			result = row.ID
			break
		}
	}
	return result
}

func FindURLID(id string) (result string) {
	result = ""
	for _, row := range urls {
		if row.ID == id {
			result = row.URL
			break
		}
	}
	return result
}
