# Task Manager API
## Description
The Task Manager API is a RESTful service for managing user tasks. It supports user authentication with JWT, role-based access control (regular and admin users), and CRUD operations for tasks. The API is built with Go, Gin, and MongoDB.
## Introduction
This API enables users to register, log in, manage tasks, and perform administrative actions like cleaning all tasks. Regular users can manage their own tasks, while admins can manage all tasks.
Features

## User registration and login with JWT authentication
Role-based access (regular and admin)
Create, read, update, and delete tasks with attributes like name, description, status, priority, due date, and owner
Admin-only task cleanup functionality

## Installation
To set up the Task Manager API locally:
#### Clone the repository (replace with your actual repo URL)
```
git clone https://github.com/HLAWR9901/a2sv_backend_go.git
```

#### Navigate to the project directory
```
cd task4/task-manager-api
```
#### Initialize a Go module
```
go mod init task_manager_api
```
#### Install dependencies
```
go mod tidy
```
#### Install required packages
```
go get github.com/gin-gonic/gin
go get github.com/dgrijalva/jwt-go
go get github.com/joho/godotenv
go get go.mongodb.org/mongo-driver/mongo
go get golang.org/x/crypto/bcrypt
go get github.com/google/uuid
```
#### Create a .env file with:
```
MONGO_URI=mongodb://localhost:27017
SECRET=your-secure-secret-key
```
#### Build and run
```
go build
go run .
```
#### Testing
Postman
Test endpoints using Postman:

Import the provided TaskManagerAPI.postman_collection.json collection (available in the repository or generated separately).
Set up an environment with baseUrl (e.g., http://localhost:8080) and variables for regularToken, adminToken, and taskId.
Run the collection to test all endpoints in sequence.

cURL
Example to create a task:
```
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Test Task",
    "description": "A test task",
    "status": "Pending",
    "priority": "Low",
    "dueDate": "2025-07-31T00:00:00Z"
  }'
```
#### Endpoints
```
POST /register
Description: Register a new user with an email, password, and optional role.
```

Request Body:
```
{
  "email": "string",
  "password": "string",
  "role": "string" // "regular" or "admin", defaults to "regular"
}
```

Response:
```
200 OK:{
  "message": "user registered successfully"
}
```
```
400 Bad Request: Invalid input or email already exists{"error": "user already exists"}
```

```
POST /login
Description: Log in a user and receive a JWT token.
```
Request Body:
```
{
  "email": "string",
  "password": "string"
}
```
Response:
```
200 OK:{
  "message": "user logged successfully",
  "token": "string"
}
```
```
401 Unauthorized: Invalid credentials{"error": "Invalid credentials"}
```

```
POST /logout
Description: Log out a user (client-side token invalidation).
```
Headers:
```
Authorization: Bearer <token>
```
Response:
```
200 OK:{
  "message": "user logged out successfully",
  "token": "-"
}
```
```
401 Unauthorized: Invalid or missing token
```
```
POST /tasks
Description: Create a new task owned by the authenticated user.
```
Headers:
```
Authorization: Bearer <token>
```
Request Body:
```
{
  "name": "string",
  "description": "string",
  "status": "string", // "Pending", "In Progress", "Completed"
  "priority": "string", // "Low", "Medium", "High"
  "dueDate": "string" // ISO 8601 format, e.g., "2025-07-31T00:00:00Z"
}
```
Response:
```
201 Created:{
  "id": "string",
  "name": "string",
  "description": "string",
  "status": "string",
  "priority": "string",
  "dueDate": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "owner": "string"
}
```
```
400 Bad Request: Invalid input
401 Unauthorized: Invalid or missing token
```
```
GET /tasks
Description: Retrieve tasks. Regular users see only their tasks; admins see all tasks.
```
Headers:
```
Authorization: Bearer <token>
```
Response:
```
200 OK:[
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "status": "string",
    "priority": "string",
    "dueDate": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "owner": "string"
  }
]
```
```
401 Unauthorized: Invalid or missing token
500 Internal Server Error: Database error
```
```
GET /tasks/:id
Description: Retrieve a specific task by ID. Only the owner or an admin can access it.
```
Headers:
```
Authorization: Bearer <token>
```
Response:
```
200 OK:{
  "id": "string",
  "name": "string",
  "description": "string",
  "status": "string",
  "priority": "string",
  "dueDate": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "owner": "string"
}
```
```
401 Unauthorized: Invalid or missing token
403 Forbidden: User lacks permission
404 Not Found: Task not found
```
```
PUT /tasks/:id
Description: Update a task by ID. Only the owner or an admin can update it.
```
Headers:
```
Authorization: Bearer <token>
```
Request Body:
```
{
  "name": "string",
  "description": "string",
  "status": "string", // "Pending", "In Progress", "Completed"
  "priority": "string", // "Low", "Medium", "High"
  "dueDate": "string" // ISO 8601 format
}
```
Response:
```
200 OK:{
  "id": "string",
  "name": "string",
  "description": "string",
  "status": "string",
  "priority": "string",
  "dueDate": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "owner": "string"
}
```
```
400 Bad Request: Invalid input
401 Unauthorized: Invalid or missing token
403 Forbidden: User lacks permission
404 Not Found: Task not found
```
```
DELETE /tasks/:id
Description: Delete a task by ID. Only the owner or an admin can delete it.
```
Headers:
```
Authorization: Bearer <token>
```
Response:
```
204 No Content: Successful deletion
401 Unauthorized: Invalid or missing token
403 Forbidden: User lacks permission
404 Not Found: Task not found
```
```
POST /tasks/clean
Description: Delete all tasks (admin only).
```
Headers:
```
Authorization: Bearer <token>
```
Request Body:
```
{}
```
Response:
```
200 OK:{
  "message": "successful cleanup"
}
```
```
401 Unauthorized: Invalid or missing token
403 Forbidden: User is not an admin
500 Internal Server Error: Database error
```
## Notes

All endpoints except /register and /login require a JWT token in the Authorization: Bearer <token> header.
The /tasks/clean endpoint is restricted to users with the admin role.
Email comparisons are case-insensitive to prevent duplicate accounts.
Task ownership is enforced: regular users can only access their own tasks, while admins can access all tasks.
Ensure MongoDB is running and the MONGO_URI and SECRET environment variables are set.

## License
This project is licensed under the MIT License.