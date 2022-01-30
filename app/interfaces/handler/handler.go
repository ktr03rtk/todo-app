package handler

import (
	"fmt"
	"net/http"
	"text/template"
	"todo-app/usecase"
)

type Handler interface {
	Start()
}

type handler struct {
	taskUsecase usecase.TaskUsecase
}

func NewHandler(u usecase.TaskUsecase) Handler {
	return &handler{taskUsecase: u}
}

func (h *handler) Start() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}

	http.HandleFunc("/task", h.taskFindAll)

	fmt.Println("server start")
	server.ListenAndServe()
}

func (h *handler) taskFindAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskUsecase.FindAll()
	if err != nil {
		fmt.Fprintf(w, "error")

		return
	}

	files := []string{"templates/layout.html", "templates/task.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", tasks); err != nil {
		fmt.Fprintf(w, "error")
	}
}
