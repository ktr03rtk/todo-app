package handler

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"todo-app/domain/model"
	"todo-app/usecase"

	"github.com/julienschmidt/httprouter"
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
	router := httprouter.New()
	router.GET("/task", h.findAllTask)
	router.GET("/task/:id", h.findTask)

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (h *handler) findAllTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tasks, err := h.taskUsecase.FindAll()
	if err != nil {
		fmt.Fprintf(w, "error")

		return
	}

	files := []string{"templates/layout.html", "templates/task_all.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", tasks); err != nil {
		fmt.Fprintf(w, "error")
	}
}

func (h *handler) findTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		fmt.Fprintf(w, "error")

		return
	}

	files := []string{"templates/layout.html", "templates/task_detail.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", task); err != nil {
		fmt.Fprintf(w, "error")
	}
}
