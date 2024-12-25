// helpers.go
package utils

import (
    "github.com/gin-gonic/gin"
)

// jsonErrorResponse handles JSON error responses
func JsonErrorResponse(c *gin.Context, status int, message string) {
    c.JSON(status, gin.H{"error": message})
}