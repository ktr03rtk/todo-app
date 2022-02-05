package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"todo-app/domain/model"

	"github.com/julienschmidt/httprouter"
)

const timeLayout = "2006-01-02"

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format(timeLayout)
	},
}

func (h *handler) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	generateHTML(w, r, nil, "layout", "home")
}

func (h *handler) findAllTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	tasks, err := h.taskUsecase.FindAll()
	if err != nil {
		errorResponse(w, r, err)
	}

	generateHTML(w, r, tasks, "layout", "task_all")
}

func (h *handler) newTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	generateHTML(w, r, nil, "layout", "task_new")
}

func (h *handler) createTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, r, err)
	}

	if err := h.taskUsecase.Create(r.PostFormValue("name"), r.PostFormValue("detail"), deadline); err != nil {
		errorResponse(w, r, err)
	}

	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *handler) findTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, r, err)
	}

	generateHTML(w, r, task, "layout", "task_detail")
}

func (h *handler) editTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	id := model.TaskID(ps.ByName("id"))

	task, err := h.taskUsecase.FindByID(id)
	if err != nil {
		errorResponse(w, r, err)
	}

	files := []string{"templates/layout.html", "templates/task_edit.html"}
	templates := template.Must(template.New("editTask").Funcs(funcMap).ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", task); err != nil {
		errorResponse(w, r, err)
	}
}

func (h *handler) updateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)
	}

	id := model.TaskID(ps.ByName("id"))

	status, err := strconv.Atoi(r.PostFormValue("status"))
	if err != nil {
		errorResponse(w, r, err)
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, r, err)
	}

	if err := h.taskUsecase.Update(id, r.PostFormValue("name"), r.PostFormValue("detail"), model.Status(status), deadline); err != nil {
		errorResponse(w, r, err)
	}

	url := fmt.Sprint("/tasks/show/", id)
	http.Redirect(w, r, url, http.StatusFound)
}
