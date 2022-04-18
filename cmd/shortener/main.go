
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

func writeURL(id string, url string){
	item := urlItem{ID: id, URL: url}
	urls = append(urls, item)
}

func findURL(id string) (result string){
	result = ""
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

func isValidURL(toTest string) bool {
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


func GetR2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:		
		id := strings.Split(string(r.URL.Path),"/")[1]
		if id == "" {
			http.Error(w, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		fmt.Println(findURL(id))
		urlFind:=findURL(id)
		if urlFind=="" {
			http.Error(w, "url not exist", http.StatusBadRequest)
			return
		}else{
			w.Header().Set("Location", urlFind)
			// 307-http.StatusTemporaryRedirect  200 -http.StatusOK
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
		

	case http.MethodPost:
		fmt.Println("post")
		link, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		l:=strings.ReplaceAll(string(link), "'", "")
		if isValidURL(l){

			id:=RandStringBytes(7)
			writeURL(id,l)
			w.WriteHeader(http.StatusCreated)
			b:=[]byte(`http://localhost:8080/`+id)
			w.Write(b)

		} else {
			fmt.Fprintf(w, "url is not valid "+l)
			http.Error(w, "url is not valid ", http.StatusBadRequest)
			//w.WriteHeader(http.StatusBadRequest)  //400
		}
		
	default:
		fmt.Fprintf(w, " only GET and POST methods are supported")
		http.Error(w, "only GET and POST methods are supported", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)  //400
	}
}


func main() {

	writeURL("12","https://pkg.go.dev/net/http")

	http.HandleFunc("/",GetR2)

    server := &http.Server{
        Addr: "localhost:8080",
    }
    log.Fatal(server.ListenAndServe())
} 