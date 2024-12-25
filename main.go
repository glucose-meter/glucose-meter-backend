package main

import (
    "context"
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "glucose-meter-backend/controllers"
    "glucose-meter-backend/database"
    "os/signal"
    "syscall"
    "time"
    // "github.com/dgrijalva/jwt-go" // For JWT handling
    // "golang.org/x/crypto/bcrypt"  // For password hashing
)

func init() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

func initDB() {
    // Initialize the connection pool
    database.InitializeDB()
    log.Println("Connected to PostgreSQL")
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

func passCtxToByPassNotUsedError(ctx context.Context) {
    // Do nothing
}

func main() {
    // Initialize the database connection pool
    initDB()

    // Set up Gin router
    router := gin.Default()
    router.Use(CORSMiddleware()) // Apply CORS middleware globally

    // Routes for Device Management 
    router.POST("/glucose/add", controllers.AddData)
    router.GET("/glucose/download", controllers.DownloadData)

    // Graceful shutdown
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        port := os.Getenv("PORT")
        if port == "" {
            port = "8080" // Fallback port if not set
        }
        if err := router.Run(":" + port); err != nil {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()

    // Wait for shutdown signal
    <-shutdown
    log.Println("Shutting down server...")

    // Graceful shutdown process
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    passCtxToByPassNotUsedError(ctx)
    defer cancel()

    // Close database connection pool
    database.CloseDB()

    log.Println("Server stopped.")
}