package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"todo-app/config"
	"todo-app/domain/model"
	"todo-app/infrastructure/persistence"

	"gorm.io/gorm"
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

func taskHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/layout.html", "templates/task.html"}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", nil)
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}

	conn := config.NewDBConn()

	id := taskCreate(conn)
	taskRead(conn, id)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/tasks", taskHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}

func taskCreate(conn *gorm.DB) model.TaskID {
	id := model.TaskID(model.CreateUUID())
	name := "test task"
	detail := "create test task"
	deadline := time.Now().Add(48 * time.Hour)

	task, err := model.CreateTask(id, name, detail, deadline)
	fmt.Printf("--------------- %+v\n", task)
	if err != nil {
		panic(err)
	}

	tp := persistence.NewTaskPersistence(conn)
	if err := tp.Create(task); err != nil {
		panic(err)
	}

	return id
}

func taskRead(conn *gorm.DB, id model.TaskID) {
	tp := persistence.NewTaskPersistence(conn)
	task, err := tp.FindByID(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("--------------- %+v\n", task)
}
