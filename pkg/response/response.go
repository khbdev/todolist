package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, gin.H{
        "success": true,
        "data":    data,
    })
}

func Error(c *gin.Context, message string, code int) {
    c.JSON(code, gin.H{
        "success": false,
        "error":   message,
    })
}