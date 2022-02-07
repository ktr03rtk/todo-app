package handler

import (
	"context"
	"log"
	"net/http"
	"time"
	"todo-app/usecase"

	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	Start()
	Stop()
}

type handler struct {
	taskUsecase    usecase.TaskUsecase
	userUsecase    usecase.UserUsecase
	sessionUsecase usecase.SessionUsecase
	server         *http.Server
}

func NewHandler(tu usecase.TaskUsecase, uu usecase.UserUsecase, su usecase.SessionUsecase) Handler {
	h := &handler{
		taskUsecase:    tu,
		userUsecase:    uu,
		sessionUsecase: su,
	}

	h.setupServer()

	return h
}

func (h *handler) Start() {
	if err := h.server.ListenAndServe(); err != nil {
		log.Fatalln("Server closed with error:", err)
	}
}

func (h *handler) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Println("Failed to gracefully shutdown:", err)
	}

	log.Println("Server shutdown")
}

func (h *handler) setupServer() {
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

	h.server = &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
}
