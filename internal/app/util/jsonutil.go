package util

import (
	"encoding/json"
	"fmt"
)

type strUrl struct {
	Url string `json:"url"`
}

type strUrlres struct {
	Url string `json:"result"`
}

func StrtoJson(u string) []byte {
	val := strUrlres{
		Url: u,
	}
	ju, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return ju
}

func JsontoUrl(b []byte) string {
	val := strUrl{}
	if err := json.Unmarshal([]byte(b), &val); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", val)
	return string(string(val.Url))
}
