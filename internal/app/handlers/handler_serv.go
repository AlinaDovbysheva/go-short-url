package handlers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/storage"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type HandlerServer struct {
	Chi *chi.Mux
	s   storage.DBurl
}

func NewHandlerServer(st storage.DBurl) *HandlerServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CookieHandle)

	h := HandlerServer{
		Chi: r,
		s:   st,
	}

	h.Chi.Get("/api/user/urls", h.HandlerServerGetUrls)
	h.Chi.Get("/{id}", h.HandlerServerGet)
	h.Chi.Post("/api/shorten", h.HandlerServerPostJSON)
	h.Chi.Post("/", h.HandlerServerPost)
	h.Chi.Post("/api/shorten/batch", h.HandlerServerPostJSONArray)
	h.Chi.Delete("/api/user/urls", h.HandlerServerPostDelArray)
	return &h
}

func CookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		nc := NewCookie()

		if err != nil {
			r.AddCookie(nc)
			http.SetCookie(w, nc)
			fmt.Println("New cookie.Value : ", nc.Value)
		}
		next.ServeHTTP(w, r)
	})
}

func NewCookie() *http.Cookie {
	nc := util.NewCookie()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookieNew := &http.Cookie{Name: "token", Value: nc, Expires: expiration}
	return cookieNew
}

func (h *HandlerServer) HandlerServerGetUrls(w http.ResponseWriter, r *http.Request) {
	//set Cookie
	cookie, _ := r.Cookie("token")
	fmt.Println("cookie.Value : ", cookie.Value)

	urlsFind, err := h.s.GetAllURLUid(cookie.Value) // storage.FindURL(id)
	fmt.Println("all url for user : ", string(urlsFind))
	if err != nil {
		w.WriteHeader(http.StatusNoContent) //204
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	w.Write(urlsFind)
}

func (h *HandlerServer) HandlerServerGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
	}
	if id == "ping" {
		err := h.s.PingDB()
		if err != nil {
			http.Error(w, "internal Server Error ", http.StatusBadRequest)
			w.WriteHeader(http.StatusInternalServerError) //500
			return
		}
		w.WriteHeader(http.StatusOK) //200
		return
	}
	// set Cookie
	cookie, _ := r.Cookie("token")
	fmt.Println(" id : ", id)
	fmt.Println(" cookie.Value : ", cookie.Value)

	urlFind, err := h.s.GetURL(id, cookie.Value) // storage.FindURL(id)
	fmt.Println(urlFind)
	if errors.Is(err, util.ErrHandler410) {
		w.WriteHeader(util.StorageErrToStatus(util.ErrHandler410)) //410
		return
	}
	if errors.Is(err, util.ErrHandler400) {
		w.WriteHeader(util.StorageErrToStatus(util.ErrHandler400))
		return
	}

	fmt.Printf("token: %v\n", cookie.Value)
	w.Header().Set("Location", urlFind)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer ?????????? ???????????????? ???? gzip-????????????, ?????????????? ?????????? ?? ????????
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ??????????????????, ?????? ???????????? ???????????????????????? gzip-????????????
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// ???????? gzip ???? ????????????????????????????, ???????????????? ????????????????????
			next.ServeHTTP(w, r)
			return
		}
		// ?????????????? gzip.Writer ???????????? ???????????????? w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// ???????????????? ?????????????????????? ???????????????? ???????????????????? ???????? gzipWriter ?????? ???????????? ????????????
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (h *HandlerServer) HandlerServerPostJSON(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}
	// set Cookie
	cookie, _ := r.Cookie("token")
	fmt.Println(" cookie.Value : ", cookie.Value)

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	u := util.JsontoURL(body)
	fmt.Println(" url:" + u)
	if util.IsValidURL(u) {
		_, jsonURL, err := h.s.PutURL(u, cookie.Value)
		fmt.Println(string(jsonURL))
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, util.ErrHandler409) {
			w.WriteHeader(util.StorageErrToStatus(util.ErrHandler409)) //409
		} else {
			w.WriteHeader(http.StatusCreated) //201
		}
		w.Write(jsonURL)
		return
	}
	fmt.Fprintf(w, "url is not valid "+u)
}

func (h *HandlerServer) HandlerServerPostJSONArray(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), 500)
		//return util.ErrHandler500
	}
	fmt.Println(" body ", string(body))

	cookie, _ := r.Cookie("token")
	fmt.Println(" cookie.Value : ", cookie.Value)

	jsonURL, err := h.s.PutURLArray(body, cookie.Value)
	if err != nil {
		fmt.Fprintf(w, "can't make short url "+string(body))
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if jsonURL != nil {
		fmt.Println(string(jsonURL))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) //201
		w.Write(jsonURL)
		return
	}
	fmt.Fprintf(w, "can't make short url "+string(body))
	w.WriteHeader(http.StatusBadRequest) //400
}

func (h *HandlerServer) HandlerServerPostDelArray(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Println(" body ", string(body))

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println(" cookie not exist : ", err)
	} else {
		fmt.Println(" cookie.Value : ", cookie.Value)
	}

	go h.s.DelURLArray(body, cookie.Value)
	w.WriteHeader(http.StatusAccepted) //202
}

func (h *HandlerServer) HandlerServerPost(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), 500)
		//return util.ErrHandler500
	}
	// set Cookie
	cookie, _ := r.Cookie("token")
	fmt.Println(" cookie.Value : ", cookie.Value)

	l := strings.ReplaceAll(string(body), "'", "")
	if util.IsValidURL(l) {
		id, _, err := h.s.PutURL(l, cookie.Value) //
		if errors.Is(err, util.ErrHandler409) {
			w.WriteHeader(util.StorageErrToStatus(util.ErrHandler409)) //409
		} else {
			w.WriteHeader(http.StatusCreated) //201
		}
		b := []byte(app.BaseURL + `/` + id)
		fmt.Println(string(b))
		w.Write(b)
		return
	}
	http.Error(w, "Url is not valid ", http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest) //400
}
