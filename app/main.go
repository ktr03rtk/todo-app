package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"todo-app/config"
	"todo-app/domain/model"
	"todo-app/infrastructure/persistence"
	"todo-app/usecase"

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

	taskCreate(conn)
	taskRead(conn)
	taskUpdate(conn)
	taskRead(conn)
	taskReadAll(conn)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/task", taskHandler)

	fmt.Println("server start")

	server.ListenAndServe()
}

func taskCreate(conn *gorm.DB) {
	name := "test task"
	detail := "create test task"
	deadline := time.Now().Add(48 * time.Hour)

	tp := persistence.NewTaskPersistence(conn)
	usecase := usecase.NewTaskUsecase(tp)
	if err := usecase.Create(name, detail, deadline); err != nil {
		panic(err)
	}
}

func taskRead(conn *gorm.DB) {
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	tp := persistence.NewTaskPersistence(conn)
	usecase := usecase.NewTaskUsecase(tp)

	task, err := usecase.FindByID(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--------------- %+v\n", task)
}

func taskUpdate(conn *gorm.DB) {
	id := model.TaskID("19742914-f296-4855-aa8d-f099727e288f")
	name := "updated task"
	detail := "updated test task"
	status := model.Completed
	deadline := time.Now().Add(48 * time.Hour)

	tp := persistence.NewTaskPersistence(conn)
	usecase := usecase.NewTaskUsecase(tp)

	err := usecase.Update(id, name, detail, status, deadline)
	if err != nil {
		panic(err)
	}
}

func taskReadAll(conn *gorm.DB) {
	tp := persistence.NewTaskPersistence(conn)
	usecase := usecase.NewTaskUsecase(tp)

	tasks, err := usecase.FindAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("find all -------------------")
	for i, t := range tasks {
		fmt.Printf("---------------%v, %#v\n", i, t)
	}
}
