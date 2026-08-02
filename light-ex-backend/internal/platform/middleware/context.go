package middleware

import "github.com/gin-gonic/gin"

// UserID returns the authenticated user's ID from Gin context.
func UserID(c *gin.Context) string {

	value, exists := c.Get(ContextUserID)

	if !exists {
		return ""
	}

	id, ok := value.(string)

	if !ok {
		return ""
	}

	return id
}