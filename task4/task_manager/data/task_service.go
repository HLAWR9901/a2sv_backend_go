package data

import (
	"errors"
	"sync"
	"task_manager/models"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("task not found")
var ErrInvalid = errors.New("invalid required field")
type TaskRepo interface{
	Create(models.Task) error
	UpdateById(models.Task)error
	DeleteByID(id uuid.UUID)error
	GetAll()([]models.Task,error)
	GetById(id uuid.UUID)(*models.Task,error)
}

type InMemoryTaskRepo struct{
	mu sync.RWMutex
	tasks []models.Task
}

func NewInMemoryRepo()*InMemoryTaskRepo{
	return &InMemoryTaskRepo{tasks:make([]models.Task,0)}
}

func (m *InMemoryTaskRepo) Create(t models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.BaseModel.ID == uuid.Nil||
		t.Name == ""||
		t.DueDate.IsZero()||
		t.CreatedAt.IsZero()||
		t.Status!=models.State(models.Pending) {
		return ErrInvalid
	}
	m.tasks = append(m.tasks, t)
	return nil
}

func (m *InMemoryTaskRepo) UpdateById(t models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i:= range m.tasks{
		if m.tasks[i].BaseModel.ID==t.BaseModel.ID{
			if t.Name!=""{
				m.tasks[i].Name=t.Name
			}
			if t.Description!=""{
				m.tasks[i].Description = t.Description
			}
			if t.Status == models.State(models.Pending) || 
				t.Status == models.State(models.InProgress) ||
				t.Status == models.State(models.Completed){
					m.tasks[i].Status = t.Status
			}
			if t.Priority == models.Importance(models.High)||
				t.Priority == models.Importance(models.Medium)||
				t.Priority == models.Importance(models.Low){
					m.tasks[i].Priority = t.Priority
			}
			if !t.DueDate.IsZero(){
				m.tasks[i].DueDate = t.DueDate
			}
			if !t.UpdatedAt.IsZero(){
				m.tasks[i].UpdatedAt = t.UpdatedAt
			}
			if t.CompletedAt!=nil{
				m.tasks[i].CompletedAt = t.CompletedAt
			}
			return nil
		}
	}
	return ErrNotFound
}

func (m *InMemoryTaskRepo) DeleteByID(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
    for i := range m.tasks {
        if m.tasks[i].ID == id {
            m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
            return nil
        }
    }
    return ErrNotFound
}

func (m *InMemoryTaskRepo) GetAll()([]models.Task,error) {
    m.mu.RLock()
    defer m.mu.RUnlock()	
    return append([]models.Task(nil), m.tasks...), nil
}

func (m *InMemoryTaskRepo) GetById(id uuid.UUID) (*models.Task,error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for i:= range m.tasks{
		if m.tasks[i].BaseModel.ID == id{
			return &m.tasks[i],nil
		}
	}
	return nil,ErrNotFound
}