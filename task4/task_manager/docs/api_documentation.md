# 📌 Task Manager API

## 📝 Description
The **Task Manager API** is a RESTful service designed to help users manage their tasks efficiently. It provides endpoints for creating, reading, updating, and deleting tasks, offering a seamless experience for task management.

---

## 📚 Introduction

The Task Manager API allows users to interact with a task management system through a set of well-defined API endpoints. Whether you're building a frontend application or integrating with other services, this API provides the necessary tools to handle tasks effectively.

---

## 🚀 Features

- ✅ Create new tasks with details like name, description, status, priority, and due date.
- 📋 Retrieve a list of all tasks or the details of a specific task.
- ✏️ Update existing tasks to change their attributes like status, priority, or content.
- 🗑️ Delete tasks that are no longer needed.

---

## ⚙️ Installation

To set up the Task Manager API on your local machine, follow these steps:

```bash
# Clone the repository
git clone https://github.com/HLAWR9901/a2sv_backend_go.git

# Navigate to the project directory
cd a2sv_backend_go/task4/task_manager_api/

# Initialize a Go module (if not already initialized)
go mod init task_manager_api

# Install dependencies
go mod tidy

# Optionally install Gin if not resolved
go get github.com/gin-gonic/gin

# Build and Run the project
go build
# or run directly
go run main.go
```
---

## 🧪 Testing

### 🔹 Postman
You can test all the API endpoints using Postman:

1. Open the [Postman workspace](https://www.postman.com/apeiron1/workspace/my-workspace/run/46807462-b0e308bf-8ee1-4913-839f-f91046e9beff)
2. Import or open the collection.
3. Run individual requests or use the **Collection Runner** for automation.

### 🔹 cURL
You can also use cURL commands to test endpoints from the terminal. Example:

```bash
curl -X POST http://localhost:3000/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Task",
    "description": "A test task entry",
    "status": "Pending",
    "priority": "Low",
    "duedate": "2025-07-31T00:00:00Z"
  }'
