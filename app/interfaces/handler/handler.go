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
	"github.com/pkg/errors"
)

type Handler interface {
	Start()
}

type handler struct {
	taskUsecase    usecase.TaskUsecase
	userUsecase    usecase.UserUsecase
	sessionUsecase usecase.SessionUsecase
}

func NewHandler(tu usecase.TaskUsecase, uu usecase.UserUsecase, su usecase.SessionUsecase) Handler {
	return &handler{
		taskUsecase:    tu,
		userUsecase:    uu,
		sessionUsecase: su,
	}
}

const timeLayout = "2006-01-02"

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format(timeLayout)
	},
}

func (h *handler) Start() {
	router := httprouter.New()

	// INFO: avoid conflict https://github.com/julienschmidt/httprouter/issues/73
	router.GET("/", h.home)
	router.GET("/tasks", h.findAllTask)
	router.GET("/tasks/new", h.newTask)
	router.POST("/tasks", h.createTask)
	router.GET("/tasks/show/:id", h.findTask)
	router.GET("/tasks/show/:id/edit", h.editTask)
	router.POST("/tasks/show/:id", h.updateTask)

	router.GET("/signup", h.signup)
	router.POST("/signup", h.signupUser)

	router.GET("/login", h.login)
	router.POST("/login", h.authenticate)
	router.GET("/logout", h.logout)

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (h *handler) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	files := []string{"templates/layout.html", "templates/home.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) findAllTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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

func (h *handler) newTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	files := []string{"templates/layout.html", "templates/task_new.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) createTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if err := r.ParseForm(); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	deadline, err := time.Parse(timeLayout, r.PostFormValue("deadline"))
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	if err := h.taskUsecase.Create(r.PostFormValue("name"), r.PostFormValue("detail"), deadline); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *handler) findTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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

	url := fmt.Sprint("/tasks/show/", id)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	files := []string{"templates/layout.html", "templates/signup.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) signupUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	if err := h.userUsecase.Signup(r.PostFormValue("email"), r.PostFormValue("password")); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	files := []string{"templates/layout.html", "templates/login.html"}
	templates := template.Must(template.ParseFiles(files...))

	if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}
}

func (h *handler) authenticate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	id, err := h.userUsecase.Authenticate(r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	session, err := h.sessionUsecase.CreateSession(id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	cookie := &http.Cookie{
		Name:     "todo_cookie",
		Value:    string(session.ID),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if err := h.sessionUsecase.DeleteSession(s.UserID); err != nil {
		errorResponse(w, err, http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *handler) session(r *http.Request) (*usecase.Session, error) {
	cookie, err := r.Cookie("todo_cookie")
	if err != nil {
		return nil, nil
	}

	session, err := h.sessionUsecase.Verify(usecase.SessionID(cookie.Value))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to verify cookie")
	}

	return session, nil
}

func errorResponse(w http.ResponseWriter, err error, errorCode int) {
	fmt.Println(err)
	http.Error(w, err.Error(), errorCode)
}
