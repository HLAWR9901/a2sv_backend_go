package main

import (
	"context"
	"log"
	"task_manager/delivery/controllers"
	"task_manager/delivery/routers"
	"task_manager/infrastructure/config"
	"task_manager/infrastructure/db"
	"task_manager/infrastructure/auth"
	"task_manager/infrastructure/passwordservice"
	"task_manager/repositories"
	"task_manager/usecases"
)

func main() {
	// Loading config - .env variables
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	// Database connection
	client, err := db.Connect(context.Background(), cfg.MongodbURI, config.DBConnectTimeOut)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	// password and auth service init
	pwSvc := passwordservice.NewBcryptService(passwordservice.DefaultCost)
	jwtSvc := auth.NewJWTService(cfg.JwtSecret, config.TokenTTL)

	// Repo init
	repo := repositories.NewRepository(client, cfg.MongodbName)
	userRepo := repositories.NewUserRepository(repo)
	taskRepo := repositories.NewTaskRepository(repo)

	// Usecase init
	userUC := usecases.NewUserUsecase(userRepo, pwSvc, jwtSvc)
	taskUC := usecases.NewTaskUsecase(taskRepo)

	// Controller set on the repos
	handler := controllers.NewHandler(userUC, taskUC)

	// Setup the router - need Port, Handler, auth service
	routers.SetupRouter(cfg.ServerPort, handler, jwtSvc)
}