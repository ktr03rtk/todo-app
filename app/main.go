package main

import (
	"todo-app/config"
	"todo-app/infrastructure/persistence"
	"todo-app/interfaces/handler"
	"todo-app/usecase"
)

func main() {
	conn := config.NewDBConn()
	taskRepository := persistence.NewTaskPersistence(conn)
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	handler := handler.NewHandler(taskUsecase)

	handler.Start()
}
