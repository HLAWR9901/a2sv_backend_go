package main

import (
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"

	"github.com/gin-gonic/gin"
)

func main() {
	//Initialize Inmemory repository
	repo := data.NewInMemoryRepo()

	r:= gin.Default()

	service := &controllers.Service{Repo: repo}
	router.SetRouter(r,service)

	r.Run(":3000")	
}