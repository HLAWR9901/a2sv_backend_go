package controllers

import (
	"net/http"
	"strings"
	"task_manager/data"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
type TaskHandler interface{
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	GetByID(c *gin.Context)
}

type Service struct{
	Repo data.TaskRepo
}

func (s *Service) Create(c *gin.Context) {
	var task models.Task
	if parseErr:=c.BindJSON(&task); parseErr!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":parseErr.Error()})
		return
	}
	task.BaseModel.ID = uuid.New()
	task.BaseModel.CreatedAt = time.Now()
	task.BaseModel.UpdatedAt = time.Now()
	//on integration to handle case for name and description
	 task.Name = strings.ToLower(task.Name)
	 task.Description = strings.ToLower(task.Description)
	if dataErr:=s.Repo.Create(task);dataErr!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":dataErr.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated,task)
}

func (s *Service) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var task models.Task
	if parseErr:=c.BindJSON(&task); parseErr!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":parseErr.Error()})
		return
	}
	task.BaseModel.ID=id
	task.BaseModel.UpdatedAt = time.Now()

	if dataErr:=s.Repo.UpdateById(task);dataErr!=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"error":dataErr.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK,gin.H{"message":"task updated"})
}

func (s *Service) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if dataErr:=s.Repo.DeleteByID(id); dataErr!=nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"error":dataErr.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK,gin.H{"message":"task deleted"})
}

func (s *Service) Get(c *gin.Context) {
	tasks,err:=s.Repo.GetAll()
	if err!=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"error":err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK,tasks)
}

func (s *Service) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	task,dataErr := s.Repo.GetById(id)
	if dataErr!=nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"error":dataErr.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK,task)
}