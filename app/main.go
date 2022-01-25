package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"todo-app/domain/model/task_model"

	"gorm.io/driver/mysql"
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

	taskCreate()
	taskRead()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/tasks", taskHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}

var Db *gorm.DB

func taskCreate() {
	var err error
	dsn := "root:password@tcp(db:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	id := task_model.TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	name := "test task"
	detail := "create test task"
	deadline := time.Now().Add(48 * time.Hour)

	task, err := task_model.CreateTask(id, name, detail, deadline)
	fmt.Printf("--------------- %+v\n", task)
	if err != nil {
		panic(err)
	}

	result := Db.Create(&task)
	fmt.Printf("--------------- %+v\n", result)
}

func taskRead() {
	id := task_model.TaskID("72c24944-f532-4c5d-a695-70fa3e72f3ab")
	var task task_model.Task
	Db.Take(&task, id)
	fmt.Printf("--------------- %+v\n", task)
}
