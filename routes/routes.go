package routes

import (
    "github.com/gin-gonic/gin"
    "glucose-meter-backend/controllers"
)

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

// SetupRouter initializes all routes and applies middleware
func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.Use(CORSMiddleware()) // Apply CORS middleware globally

    // Device routes
    deviceRoutes := r.Group("/glucose")
    {
        deviceRoutes.POST("/add", controllers.AddData)
        deviceRoutes.GET("/download", controllers.DownloadData)
    }

    return r
}
