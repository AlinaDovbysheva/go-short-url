package handlers

import (
	"bytes"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			appH := NewHandlerServer()
			appH.ServeHTTP(w, rPost)
			resp := w.Result()

			assert.Equal(t, tt.codepost1, resp.StatusCode)
			rGetIDjson, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			rGetID := util.JsontoURLRes(rGetIDjson) //{"result":"<shorten_url>"}
			rGet := httptest.NewRequest(http.MethodGet, string(rGetID), nil)
			w = httptest.NewRecorder()
			appH.ServeHTTP(w, rGet)
			resp = w.Result()

			rGetURL := w.Header().Get("Location")
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.codeget, resp.StatusCode)
			assert.Equal(t, tt.value, string(rGetURL))

		})
	}
}
