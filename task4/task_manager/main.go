package main

import (
    "context"
    "log"
    "os"
    "task_manager/controllers"
    "task_manager/data"
    "task_manager/router"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    Init()

    uri := os.Getenv("URI")
    if uri == "" {
        log.Fatal("Environment variable URI is not set")
    }

    conn := connect(uri)
    defer func() {
        if err := conn.Disconnect(context.Background()); err != nil {
            log.Printf("Error disconnecting from MongoDB: %v", err)
        }
    }()

    repo := data.NewRepo(conn, false)
    handler := controllers.SetHandler(repo)
    r := router.NewRouter(handler)

    if err := r.Run(":3000"); err != nil {
        log.Fatal(err)
    }
}

func Init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}

func connect(uri string) *mongo.Client {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)
    for attempt := 1; attempt <= 3; attempt++ {
        client, err := mongo.Connect(ctx, clientOpts)
        if err == nil {
            if err := client.Ping(ctx, nil); err == nil {
                return client
            }
            log.Printf("Ping failed on attempt %d: %v", attempt, err)
        } else {
            log.Printf("Connection failed on attempt %d: %v", attempt, err)
        }
        time.Sleep(time.Second * time.Duration(attempt))
    }
    log.Fatal("Failed to connect to MongoDB after retries")
    return nil
}