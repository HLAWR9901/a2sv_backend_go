package controllers

import (
    "net/http"
    "task_manager/data"
    "task_manager/models"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type Handler struct {
    Repo *data.Repo
}

func SetHandler(repo *data.Repo) *Handler {
    return &Handler{Repo: repo}
}

func (h *Handler) Create(c *gin.Context) {
    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task.ID = uuid.New()
    now := time.Now()
    task.CreatedAt = now
    task.UpdatedAt = now
    task.Status = models.Pending
    if err := task.Validate(); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.Repo.Create(&task); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusCreated, task)
}

func (h *Handler) Update(c *gin.Context) {
    id := c.Param("id")
    if _, err := uuid.Parse(id); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    task.ID = uuid.MustParse(id)
    task.UpdatedAt = time.Now()
    if err := h.Repo.Update(id, task); err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, task)
}

func (h *Handler) Delete(c *gin.Context) {
    id := c.Param("id")
    if _, err := uuid.Parse(id); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    if err := h.Repo.Delete(id); err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) GetById(c *gin.Context) {
    id := c.Param("id")
    if _, err := uuid.Parse(id); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }
    task, err := h.Repo.GetById(id)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, task)
}

func (h *Handler) GetAll(c *gin.Context) {
    tasks, err := h.Repo.GetAll()
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, tasks)
}