package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewCookie() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ""
	}
	return hex.EncodeToString(b)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

var (
	ErrHandler404    = errors.New("not found")
	ErrHandler400    = errors.New("bad request")
	ErrHandler500    = errors.New("internal server error")
	ErrAlreadyExists = errors.New("already exists")
)

func storageErrToStatus(err error) int {
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrHandler500:
		return http.StatusInternalServerError
	case ErrHandler404:
		return http.StatusNotFound
	case ErrHandler400:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
