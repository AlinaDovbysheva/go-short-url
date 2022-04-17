
package main

import (
	"math/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"log"
	"strings"
)

type urlItem struct {
    ID string 
    URL string 
}

var urls []urlItem

func writeUrl(id string, url string){
	item := urlItem{ID: id, URL: url}
	urls = append(urls, item)
}

func findUrl(id string) (result string){
	result = "not found"
	for _, row := range urls{
		if row.ID == id {
			result = row.URL
			break
		}
	}
	return result
 }
 
 const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

 func RandStringBytes(n int) string {
	 b := make([]byte, n)
	 for i := range b {
		 b[i] = letterBytes[rand.Intn(len(letterBytes))]
	 }
	 return string(b)
 }

func isValidUrl(toTest string) bool {
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

func isUrl(str string) (bool, error) {
    parsedUrl, err := url.Parse(str)
    return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "", err
}

func GetR2(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:		
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		fmt.Println(findUrl(id))
		w.Header().Set("Location", findUrl(id))
		// устанавливаем статус-код 307-http.StatusTemporaryRedirect  200 -http.StatusOK
		w.WriteHeader(http.StatusTemporaryRedirect)
	

	case http.MethodPost:
		fmt.Println("post")
		link, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		l:=strings.ReplaceAll(string(link), "'", "")
		if isValidUrl(l){

			id:=RandStringBytes(7)
			writeUrl(id,l)
			w.WriteHeader(http.StatusCreated)
			b:=[]byte(`localhost:8080/`+id)
			w.Write(b)

		} else {
			fmt.Fprintf(w, "url is not valid "+l)
			//http.Error(w, "url is not valid", http.StatusMethodNotAllowed)
		}
		
	default:
		fmt.Fprintf(w, " only GET and POST methods are supported")
		http.Error(w, "only GET and POST methods are supported", http.StatusMethodNotAllowed)
	}
}


func main() {

	http.HandleFunc("/",GetR2)

    server := &http.Server{
        Addr: "localhost:8080",
    }
    log.Fatal(server.ListenAndServe())
} 