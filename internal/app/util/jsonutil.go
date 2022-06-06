package util

import (
	"encoding/json"
)

type strURL struct {
	URL string `json:"url"`
}

type strURLres struct {
	URL string `json:"result"`
}

func StrtoJSON(u string) []byte {
	val := strURLres{
		URL: u,
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
	return string(string(val.URL))
}

func JsontoURLRes(b []byte) string {
	val := strURLres{}
	if err := json.Unmarshal([]byte(b), &val); err != nil {
		panic(err)
	}
	return string(string(val.URL))
}
