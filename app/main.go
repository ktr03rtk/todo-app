package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "first")
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "second")
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}
