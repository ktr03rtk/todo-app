package main

import (
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

	handler.Start()
}
