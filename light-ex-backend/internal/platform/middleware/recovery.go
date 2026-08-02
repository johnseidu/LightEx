package middleware

import (
	"net/http"
	"runtime/debug"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// Recovery recovers from panics and returns
// a standardized Internal Server Error response.
func Recovery(log *slog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {

			if recovered := recover(); recovered != nil {

				log.Error(
					"Unhandled panic recovered",
					slog.Any("panic", recovered),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("client_ip", c.ClientIP()),
					slog.String("stack_trace", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"error": "Internal Server Error",
					},
				)
			}
		}()

		c.Next()
	}
}
