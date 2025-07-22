package controllers

import (
	"net/http"
	"os"
	"time"

	"task_manager/data"
	"task_manager/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// General
type Handler struct {
	Repo *data.Repo
}

func SetHandler(repo *data.Repo) *Handler {
	return &Handler{Repo: repo}
}

// User Handlers
func (h *Handler) Register(c *gin.Context) {
	var authUser models.AuthUser
	if err := c.BindJSON(&authUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authUser.Password = string(pass)
	if h.Repo.Exists(authUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}
	user := models.User{
		ID:       uuid.NewString(),
		Email:    authUser.Email,
		Password: authUser.Password,
		Role:     authUser.Role,
	}
	if user.Role == "" {
		user.Role = models.Regular
	}
	if err := h.Repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var authUser models.AuthUser
	if err := c.BindJSON(&authUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Repo.GetUser(authUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user logged successfully", "token": tokenString})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user logged out successfully", "token": "-"})
}

// Task Operations
func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = uuid.NewString()
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	task.Status = models.Pending
	task.Owner = userID.(string)
	if err := task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	existingTask, err := h.Repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if existingTask.Owner != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this task"})
		return
	}
	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateData.Name != "" {
		existingTask.Name = updateData.Name
	}
	if updateData.Description != "" {
		existingTask.Description = updateData.Description
	}
	if updateData.Status != "" {
		existingTask.Status = updateData.Status
	}
	if updateData.Priority != "" {
		existingTask.Priority = updateData.Priority
	}
	if updateData.DueDate != nil {
		existingTask.DueDate = updateData.DueDate
	}
	existingTask.UpdatedAt = time.Now()
	if err := h.Repo.Update(id, *existingTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existingTask)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	existingTask, err := h.Repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if existingTask.Owner != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this task"})
		return
	}
	if err := h.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetById(c *gin.Context) {
	id := c.Param("id")
	task, err := h.Repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if task.Owner != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	var tasks []models.Task
	var err error
	if role == "admin" {
		tasks, err = h.Repo.GetAll()
	} else {
		tasks, err = h.Repo.GetByOwner(userID.(string))
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) CleanAll(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}
	if err := h.Repo.CleanAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful cleanup"})
}