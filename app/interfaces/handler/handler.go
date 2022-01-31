package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
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

const timeLayout = "2006-01-02"

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format(timeLayout)
	},
}

func (h *handler) Start() {
	router := httprouter.New()
	router.GET("/tasks", h.findAllTask)
	router.GET("/tasks/:id", h.findTask)
	router.GET("/tasks/:id/edit", h.editTask)

	router.POST("/tasks/:id", h.updateTask)

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (h *handler) findAllTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tasks, err := h.taskUsecase.FindAll()
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	files := []string{"templates/layout.html", "templates/task_all.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", tasks); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) findTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	files := []string{"templates/layout.html", "templates/task_detail.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", task); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) editTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	files := []string{"templates/layout.html", "templates/task_edit.html"}
	templates := template.Must(template.New("editTask").Funcs(funcMap).ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", task); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) updateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	id := model.TaskID(ps.ByName("id"))

	status, err := strconv.Atoi(r.PostFormValue("status"))
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	if err := h.taskUsecase.Update(model.TaskID(id), r.PostFormValue("name"), r.PostFormValue("detail"), model.Status(status), deadline); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	url := fmt.Sprint("/tasks/", id)
	http.Redirect(w, r, url, http.StatusFound)
}

func errorResponse(w http.ResponseWriter, err error, errorCode int) {
	fmt.Println(err)
	http.Error(w, err.Error(), errorCode)
}
