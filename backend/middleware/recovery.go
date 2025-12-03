package middleware

import "github.com/gin-gonic/gin"

// Recovery returns the default recovery middleware (wrap to customize).
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
