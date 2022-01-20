package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!!")
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "first")
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "second")
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/layout.html", "templates/task.html"}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", nil)
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/tasks", taskHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}
