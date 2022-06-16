package handlers

import (
	"bytes"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlerServer_DelArray(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		codepost1 int
	}{
		{name: "del url", value: "[\"9773005\",\"836151\",\"2806887\",\"8055539\",\"8919971\",\"9189366\",\"9498652\",\"182760\",\"8555657\",\"5489267\",\"5728872\"]", codepost1: 202},
		//{name: "simple http url", value: "https://practicum.yandex.ru/learn/go-developer/courses/9dd689b5-2524-42fb-8eef-b4e6797cbea1/sprints/21254/topics/70c00167-d0e3-4f57-a181-b5cb63c13a55/lessons/27c5eb58-ac12-4f1b-a601-5bf8c87b7ea1/", codepost1: 201, codeget: 307},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(tt.value)
			rPost := httptest.NewRequest(http.MethodDelete, "/api/user/urls", body)

			w := httptest.NewRecorder()
			db := storage.NewInMap()
			th := NewHandlerServer(db)

			nc1 := util.NewCookie()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			nc := &http.Cookie{Name: "token", Value: nc1, Expires: expiration}

			rPost.AddCookie(nc)
			http.SetCookie(w, nc)
			fmt.Println("New cookie.Value : ", nc.Value)

			appH := http.HandlerFunc(th.HandlerServerPostDelArray)
			appH.ServeHTTP(w, rPost)

			resp := w.Result()
			assert.Equal(t, tt.codepost1, resp.StatusCode)
			_ = resp.Body.Close()

		})
	}
}

func TestHandlerServer_HandlerServerGet(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		codepost1 int
	}{
		{name: "del url", value: "9773005", codepost1: 202},
		//{name: "simple http url", value: "https://practicum.yandex.ru/learn/go-developer/courses/9dd689b5-2524-42fb-8eef-b4e6797cbea1/sprints/21254/topics/70c00167-d0e3-4f57-a181-b5cb63c13a55/lessons/27c5eb58-ac12-4f1b-a601-5bf8c87b7ea1/", codepost1: 201, codeget: 307},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rPost := httptest.NewRequest(http.MethodGet, "/9773005", nil)

			w := httptest.NewRecorder()
			db := storage.NewInMap()
			th := NewHandlerServer(db)

			nc1 := util.NewCookie()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			nc := &http.Cookie{Name: "token", Value: nc1, Expires: expiration}

			rPost.AddCookie(nc)
			http.SetCookie(w, nc)
			fmt.Println("New cookie.Value : ", nc.Value)

			appH := http.HandlerFunc(th.HandlerServerGet)
			appH.ServeHTTP(w, rPost)

			//resp := w.Result()
			//assert.Equal(t, tt.codepost1, resp.StatusCode)

		})
	}
}

func TestHandlerServer_HandlerServerMain(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		codepost1 int
		codeget   int
	}{
		{name: "simple http url", value: "{\"url\":\"https://www.youtube.com/watch?v=EdufHyUAJt4\"}", codepost1: 201, codeget: 307},
		//{name: "simple http url", value: "https://practicum.yandex.ru/learn/go-developer/courses/9dd689b5-2524-42fb-8eef-b4e6797cbea1/sprints/21254/topics/70c00167-d0e3-4f57-a181-b5cb63c13a55/lessons/27c5eb58-ac12-4f1b-a601-5bf8c87b7ea1/", codepost1: 201, codeget: 307},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(tt.value)
			rPost := httptest.NewRequest(http.MethodPost, "/api/shorten", body)

			w := httptest.NewRecorder()
			db := storage.NewInMap()
			th := NewHandlerServer(db)

			nc1 := util.NewCookie()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			nc := &http.Cookie{Name: "token", Value: nc1, Expires: expiration}

			rPost.AddCookie(nc)
			http.SetCookie(w, nc)
			fmt.Println("New cookie.Value : ", nc.Value)

			appH := http.HandlerFunc(th.HandlerServerPostJSON)
			appH.ServeHTTP(w, rPost)

			resp := w.Result()

			assert.Equal(t, tt.codepost1, resp.StatusCode)
			rGetIDJson, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			rGetID := util.JsontoURLRes(rGetIDJson) //{"result":"<shorten_url>"}
			fmt.Println("rGetID: ", rGetID)

			rGet := httptest.NewRequest(http.MethodGet, string(rGetID), nil)
			w = httptest.NewRecorder()

			rGet.AddCookie(nc)
			http.SetCookie(w, nc)
			appH = http.HandlerFunc(th.HandlerServerGet)
			appH.ServeHTTP(w, rGet)
			resp = w.Result()

			rGetURL := w.Header().Get("Location")
			require.NoError(t, err)

			fmt.Println("resp: ", rGetURL)

			err = resp.Body.Close()
			require.NoError(t, err)

			//assert.Equal(t, tt.codeget, resp.StatusCode)
			//assert.Equal(t, tt.value, string(rGetURL))

		})
	}
}
