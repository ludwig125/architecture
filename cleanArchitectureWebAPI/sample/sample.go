package main

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte("Hello, Gophers!")); err != nil {
		log.Println(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	if err := http.ListenAndServe("localhost:3000", mux); err != nil {
		log.Println(err)
	}
}
