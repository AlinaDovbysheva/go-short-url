package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type HandlerServer struct {
	//chi *chi.Mux
	s storage.DBurl
}

func NewHandlerServer(st storage.DBurl) *HandlerServer {
	/*r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)*/

	h := HandlerServer{
		//chi: r,
		s: st,
	}

	/*h.chi.Get("/{id}", h.HandlerServerGet)
	h.chi.Post("/api/shorten", h.HandlerServerPostJSON)
	h.chi.Post("/", h.HandlerServerPost)*/

	return &h
}

func (h *HandlerServer) HandlerServerGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	fmt.Println(" id:" + id)
	urlFind, _ := h.s.GetURL(id) // storage.FindURL(id)
	fmt.Println(urlFind)
	if urlFind == "" {
		http.Error(w, "Url not exist ", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", urlFind)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307

}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (h *HandlerServer) HandlerServerPostJSON(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}
	body, err := io.ReadAll(reader)

	//body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	u := util.JsontoURL(body)
	fmt.Println(string(body) + " url:" + u)
	if util.IsValidURL(u) {
		id, _ := h.s.PutURL(u) //storage.WriteURL(u)
		jsonURL := util.StrtoJSON(app.BaseURL + `/` + id)
		fmt.Println(string(jsonURL))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) //201
		w.Write(jsonURL)
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
		id, _ := h.s.PutURL(l)            //storage.WriteURL(l)
		w.WriteHeader(http.StatusCreated) //201
		b := []byte(app.BaseURL + `/` + id)
		fmt.Println(string(b))
		w.Write(b)
		return
	}
	fmt.Fprintf(w, "Url is not valid "+l)
	http.Error(w, "Url is not valid ", http.StatusBadRequest)
	//w.WriteHeader(http.StatusBadRequest)  //400
}
