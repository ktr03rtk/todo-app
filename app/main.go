package main

import (
	"fmt"
	"net/http"
)

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

	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}
