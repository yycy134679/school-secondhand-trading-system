package middleware

import "github.com/gin-gonic/gin"

// Logger is a placeholder for request logging middleware.
func Logger() gin.HandlerFunc {
	return gin.Logger()
}
