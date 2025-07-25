# Task Manager API Documentation

## Overview
The Task Manager API is a secure RESTful service that allows users to manage their tasks. It provides user authentication, task creation, updating, deletion, and retrieval. Built with Go and MongoDB, it's designed for reliability and performance.

## Key Features
- 🔐 JWT-based authentication
- ✅ Create, read, update, and delete tasks
- ⚡️ Context timeouts for all operations
- 🔒 Ownership-based resource access
- 📦 MongoDB data storage

## Requirements
- Go 1.18+
- MongoDB 5.0+
- Make (optional)

## Installation & Setup
```bash
# Clone the repository
git clone https://github.com/hlawr9901/a2sv_backend_go.git
cd a2sv_backend_go/task4/task-manager/

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
```

## Configuration (`.env` file)
```env
MONGODB_URI="mongodb://localhost:27017"
MONGODB_DBNAME="task_manager"
JWT_SECRET="ordinary-daylight-losecontrol-dancing"
SERVER_PORT="8080"
```

## Running the Server
```bash
# Start the server
go run main.go

# Or using Make
make run
```

## API Endpoints

### Authentication
| Method | Endpoint    | Description          | Request Body                             |
|--------|-------------|----------------------|------------------------------------------|
| POST   | `/register` | Register new user    | `{email: string, password: string}`     |
| POST   | `/login`    | Login existing user  | `{email: string, password: string}`     |
| DELETE | `/delete`   | Delete user account  | `{password: string}`                    |
| GET    | `/logout`   | Logout user          | (Requires Authorization header)         |

### Tasks
| Method | Endpoint       | Description                | Parameters              |
|--------|----------------|----------------------------|-------------------------|
| POST   | `/task`        | Create new task            | None                    |
| PUT    | `/task/:id`    | Update existing task       | `id` in URL path        |
| DELETE | `/task/:id`    | Delete task                | `id` in URL path        |
| GET    | `/task/:id`    | Get single task by ID      | `id` in URL path        |
| GET    | `/task`        | Get all tasks for user     | None                    |

## Request Examples

### User Registration
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

### Task Creation
```http
POST /task
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "Complete API documentation",
  "description": "Write comprehensive API docs",
  "status": "pending"
}
```

### Get All Tasks
```http
GET /task
Authorization: Bearer <your_jwt_token>
```

## Authentication
Include JWT token in Authorization header for protected routes:
```
Authorization: Bearer <your_jwt_token>
```

## Testing
Import the Postman collection to test all endpoints:
1. Open Postman
2. Click "Import"
3. Select "Raw Text"
4. Paste the collection JSON provided in the project repository

## Error Responses
The API returns standardized error responses:
```json
{
  "error": "Error message here"
}
```

Common error codes:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Access to resource denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side issue

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.