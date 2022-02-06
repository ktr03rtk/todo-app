package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"todo-app/domain/model"
	"todo-app/usecase"

	"github.com/julienschmidt/httprouter"
)

const timeLayout = "2006-01-02"

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format(timeLayout)
	},
}

type data struct {
	Session *usecase.Session
	Tasks   []*model.Task
	Task    *model.Task
}

func (h *handler) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	} else {
		generateHTML(w, r, nil, "layout", "home")
	}
}

func (h *handler) findAllTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	tasks, err := h.taskUsecase.FindAll()
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	d := &data{
		Session: s,
		Tasks:   tasks,
	}

	generateHTML(w, r, d, "layout", "task_all")
}

func (h *handler) newTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		generateHTML(w, r, nil, "layout", "task_new")
	}
}

func (h *handler) createTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)

		return
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	if err := h.taskUsecase.Create(*s, r.PostFormValue("name"), r.PostFormValue("detail"), deadline); err != nil {
		errorResponse(w, r, err)

		return
	}

	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *handler) findTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	d := &data{
		Session: s,
		Task:    task,
	}

	generateHTML(w, r, d, "layout", "task_detail")
}

func (h *handler) editTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	files := []string{"/opt/templates/layout.html", "/opt/templates/task_edit.html"}
	templates := template.Must(template.New("editTask").Funcs(funcMap).ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", task); err != nil {
		errorResponse(w, r, err)
	}
}

func (h *handler) updateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)

		return
	}

	id := model.TaskID(ps.ByName("id"))

	status, err := strconv.Atoi(r.PostFormValue("status"))
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, r, err)

		return
	}

	if err := h.taskUsecase.Update(*s, id, r.PostFormValue("name"), r.PostFormValue("detail"), model.Status(status), deadline); err != nil {
		errorResponse(w, r, err)

		return
	}

	url := fmt.Sprint("/tasks/show/", id)
	http.Redirect(w, r, url, http.StatusFound)
}
