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
	h.Post("/api/shorten", h.HandlerServerPostJson)
	h.Post("/", h.HandlerServerPost)

	//storage s
	return &h
}

func (h *HandlerServer) HandlerServerGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	urlFind := storage.FindURL(id)
	fmt.Println(urlFind)
	if urlFind == "" {
		http.Error(w, "Url not exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", urlFind)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307

}

func (h *HandlerServer) HandlerServerPostJson(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := util.JsontoUrl(body)
	fmt.Println(string(body) + " url:" + u)
	if util.IsValidURL(u) {
		id := storage.WriteURL(u)
		jsonUrl := util.StrtoJson(`http://` + app.ServerURL + `/` + id)
		fmt.Println(string(jsonUrl))
		w.WriteHeader(http.StatusCreated) //201
		w.Write(jsonUrl)
		return
	}
	fmt.Fprintf(w, "url is not valid "+u)
	http.Error(w, "Url is not valid ", http.StatusBadRequest)
	//w.WriteHeader(http.StatusBadRequest)  //400
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
	}
	fmt.Fprintf(w, "Url is not valid "+l)
	http.Error(w, "Url is not valid ", http.StatusBadRequest)
	//w.WriteHeader(http.StatusBadRequest)  //400
}
