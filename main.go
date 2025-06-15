package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
		IdleTimeout:  3 * time.Second,
	}

	http.HandleFunc("/", rootPath)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func rootPath(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		log.Fatal(err)
	}
}
