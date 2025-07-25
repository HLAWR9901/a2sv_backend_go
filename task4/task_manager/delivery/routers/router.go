package routers

import (
	"task_manager/delivery/controllers"
	"task_manager/domain"
	middleware "task_manager/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter(port string, handler *controllers.Handler, jwtSvc domain.IAuthService) {
	router := gin.Default()
	// Entry routes
	router.GET("/healthz", handler.HandleHealthCheck)
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtSvc)
	router.DELETE("/delete", authMiddleware, handler.Delete)
	router.GET("/logout", authMiddleware, handler.Logout)

	// Task routes
	taskGroup := router.Group("/task")
	// Middleware on all the task routes
	taskGroup.Use(authMiddleware)
	{
		taskGroup.POST("", handler.CreateTask)
		taskGroup.PUT("/:id", handler.UpdateTask)
		taskGroup.DELETE("/:id", handler.DeleteTask)
		taskGroup.GET("/:id", handler.GetTaskByID)
		taskGroup.GET("", handler.GetTaskByUser)
	}

	router.Run(":" + port)
}