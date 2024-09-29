package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(30*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(responseTimeOut),

		// timeout.WithResponse(transform_global.ResponseTimeOut),
	)
}
func responseTimeOut(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code": http.StatusGatewayTimeout,
		"msg":  "timeout",
	})
}
