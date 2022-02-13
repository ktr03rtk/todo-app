package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo-app/config"
	"todo-app/domain/service"
	"todo-app/infrastructure/persistence"
	"todo-app/interfaces/handler"
	"todo-app/usecase"
)

func main() {
	conn := config.NewDBConn()
	taskRepository := persistence.NewTaskPersistence(conn)
	userRepository := persistence.NewUserPersistence(conn)
	sessionRepository := persistence.NewSessionPersistence(conn)
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	userService := service.NewUService(userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userService)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepository)

	handler := handler.NewHandler(taskUsecase, userUsecase, sessionUsecase)

	go func() {
		handler.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Println("Caught SIGTERM, shutting down")

	handler.Stop()
}
