package handlers

import (
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"io"
	"net/http"
	"strings"
)

type (
	HandlerServer struct {
	}
)

func NewHandlerServer() *HandlerServer {
	return &HandlerServer{}
}

func (h *HandlerServer) HandlerServerMain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandlerServerGet(w, r)
	case http.MethodPost:
		h.HandlerServerPost(w, r)
	default:
		fmt.Fprintf(w, " only GET and POST methods are supported")
		http.Error(w, "only GET and POST methods are supported", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)  //400
	}
}

func (h *HandlerServer) HandlerServerGet(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(string(r.URL.Path), "/")[1]
	if id == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	urlFind := storage.FindURLID(id)
	fmt.Println(urlFind)
	if urlFind == "" {
		http.Error(w, "Url not exist", http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Location", urlFind)
		// 307-http.StatusTemporaryRedirect
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *HandlerServer) HandlerServerPost(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	l := strings.ReplaceAll(string(link), "'", "")
	if util.IsValidURL(l) {
		id := storage.WriteURL(l)
		w.WriteHeader(http.StatusCreated) //201
		b := []byte(`http://` + app.ServerURL + `/` + id)
		fmt.Println(string(b))
		w.Write(b)
		return
	} else {
		fmt.Fprintf(w, "url is not valid "+l)
		http.Error(w, "url is not valid ", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)  //400
	}
}
