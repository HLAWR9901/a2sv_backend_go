package controllers

import (
	"net/http"
	"task_manager/domain"
	"task_manager/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email" binding:"required,email"`
}

type TaskInputDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending inprogress completed"`
}

type TaskOutputDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Handler struct {
	Userusecase *usecases.UserUsecase
	Taskusecase *usecases.TaskUsecase
}

func NewHandler(uc *usecases.UserUsecase, tc *usecases.TaskUsecase) *Handler {
	return &Handler{Userusecase: uc, Taskusecase: tc}
}

func (h *Handler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) Register(c *gin.Context) {
	// UserDTO has no password - new redundant created to recieve request body
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidPayload})
		return
	}
	// Change the json based struct to domain.User
	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	
	result, err := h.Userusecase.Create(c.Request.Context(), &user)
	if err != nil {
		switch err.Error() {
		case domain.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case domain.ErrInvalidPayload:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		}
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"data":    UserDTO{ID: result.ID, Email: result.Email},
	})
}

func (h *Handler) Login(c *gin.Context) {
	// UserDTO has no password - new redundant created to recieve request body
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidPayload})
		return
	}
	// Change the json based struct to domain.User
	user := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	
	result, token, err := h.Userusecase.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"data":  UserDTO{ID: result.ID, Email: result.Email},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *Handler) Delete(c *gin.Context) {
	// Request extract only the password for double checking authenticity
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidPayload})
		return
	}
	// Retrieve the userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	// Valid password can only delete the user
	if err := h.Userusecase.VerifyPassword(c.Request.Context(), userID.(string), req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	
	if err := h.Userusecase.DeleteUserAndTasks(c.Request.Context(), userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req TaskInputDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidPayload})
		return
	}
	// Extract the userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	// Assemble ... avengers lol
	task := domain.Task{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.State(req.Status),
		UserID:      userID.(string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.Taskusecase.Create(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		return
	}
	
	c.JSON(http.StatusCreated, TaskOutputDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

func (h *Handler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	var req TaskInputDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidPayload})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	task, err := h.Taskusecase.UpdateTask(c.Request.Context(), userID.(string), taskID, domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.State(req.Status),
	})
	if err != nil {
		switch err.Error() {
		case domain.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		}
		return
	}
	
	c.JSON(http.StatusOK, TaskOutputDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.Taskusecase.DeleteTask(c.Request.Context(), userID.(string), taskID); err != nil {
		switch err.Error() {
		case domain.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func (h *Handler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing task ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	task, err := h.Taskusecase.GetTaskByID(c.Request.Context(), userID.(string), taskID)
	if err != nil {
		switch err.Error() {
		case domain.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		}
		return
	}

	c.JSON(http.StatusOK, TaskOutputDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

func (h *Handler) GetTaskByUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	tasks, err := h.Taskusecase.GetTasksByUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternalServer})
		return
	}

	response := make([]TaskOutputDTO, len(tasks))
	for i, task := range tasks {
		response[i] = TaskOutputDTO{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			UserID:      task.UserID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}