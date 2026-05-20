package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"success": false, "error": message})
}
