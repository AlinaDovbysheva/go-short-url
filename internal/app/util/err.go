package util

import (
	"errors"
	"net/http"
)

var (
	ErrHandler404    = errors.New("not found")
	ErrHandler409    = errors.New("url exist in DB")
	ErrHandler400    = errors.New("bad request")
	ErrHandler500    = errors.New("internal server error")
	ErrAlreadyExists = errors.New("already exists")
)

func StorageErrToStatus(err error) int {
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrHandler500:
		return http.StatusInternalServerError
	case ErrHandler404:
		return http.StatusNotFound
	case ErrHandler400:
		return http.StatusBadRequest
	case ErrHandler409:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
