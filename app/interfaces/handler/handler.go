package handler

import (
	"fmt"
	"log"
	"net/http"
	"todo-app/usecase"

	"github.com/julienschmidt/httprouter"
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

	router.GET("/signup", h.signUp)
	router.POST("/signup", h.signupUser)

	router.GET("/login", h.login)
	router.POST("/login", h.authenticate)
	router.GET("/logout", h.logout)

	router.GET("/err", h.err)

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":8080", router))
}
