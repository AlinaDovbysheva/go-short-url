package handlers

import (
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"strings"
)

type HandlerServer struct {
	*chi.Mux
}

func NewHandlerServer() *HandlerServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := HandlerServer{
		Mux: r,
	}

	h.Get("/{id}", h.HandlerServerGet)
	h.Post("/", h.HandlerServerPost)

	return &h
}

func (h *HandlerServer) HandlerServerGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
		w.WriteHeader(http.StatusTemporaryRedirect) // 307
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
