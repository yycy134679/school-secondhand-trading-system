package response

import (
	"github.com/gin-gonic/gin"
)

// Error 返回错误响应
func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": message,
	})
}

// Success 返回成功响应
func Success(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}