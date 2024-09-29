package middleware

import "github.com/gin-gonic/gin"

func (m *Middleware) InternalMiddleware(c *gin.Context) {
	userID := c.GetHeader("Authorization")
	if userID != "1aeacd11-ca80-4366-b72d-d85fd89de8e9" {
		c.Status(404)
		c.Abort()
		return
	}

	c.Next()
}
