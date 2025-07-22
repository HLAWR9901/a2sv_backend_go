package router

import (
	"task_manager/controllers"
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine, handler *controllers.Handler, AuthCheck gin.HandlerFunc) {
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/logout", AuthCheck, handler.Logout)

	tasks := r.Group("/tasks", AuthCheck)
	{
		tasks.POST("", handler.Create)
		tasks.PUT("/:id", handler.Update)
		tasks.DELETE("/:id", handler.Delete)
		tasks.GET("", handler.GetAll)
		tasks.GET("/:id", handler.GetById)
		tasks.POST("/clean", handler.CleanAll)
	}
}