package router

import (
    "task_manager/controllers"
    "github.com/gin-gonic/gin"
)

func NewRouter(handler *controllers.Handler) *gin.Engine {
    router := gin.Default()
    router.RedirectTrailingSlash = false

    tasks := router.Group("/tasks")
    {
        tasks.GET("", handler.GetAll)
        tasks.GET("/:id", handler.GetById)
        tasks.POST("", handler.Create)
        tasks.PUT("/:id", handler.Update)
        tasks.DELETE("/:id", handler.Delete)
    }
    return router
}