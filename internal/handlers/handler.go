package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// HandleGet responds to GET requests.
func HandleGet(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "GET request successful"})
}

// HandlePost responds to POST requests.
func HandlePost(c *gin.Context) {
    var json map[string]interface{}
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "POST request successful", "data": json})
}