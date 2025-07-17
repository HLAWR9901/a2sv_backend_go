package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
)

func SetRouter(r *gin.Engine, s *controllers.Service){
	tasks := r.Group("/tasks")
	{
		tasks.GET("/",s.Get)
		tasks.GET("/:id",s.GetByID)
		tasks.POST("/",s.Create)
		tasks.PUT("/:id",s.Update)
		tasks.DELETE("/:id",s.Delete)
	}
}