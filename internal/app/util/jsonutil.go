package util

import (
	"encoding/json"
)

type strURL struct {
	Url string `json:"url"`
}

type strURLres struct {
	Url string `json:"result"`
}

func StrtoJson(u string) []byte {
	val := strURLres{
		Url: u,
	}
	ju, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return ju
}

func JsontoURL(b []byte) string {
	val := strURL{}
	if err := json.Unmarshal([]byte(b), &val); err != nil {
		panic(err)
	}
	return string(string(val.Url))
}

func JsontoURLRes(b []byte) string {
	val := strURLres{}
	if err := json.Unmarshal([]byte(b), &val); err != nil {
		panic(err)
	}
	return string(string(val.Url))
}
