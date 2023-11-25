package general

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/mwrooth/gateway/openapi/go"
)

func (s *ServiceGeneral) Login() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Logged in successfully!",
		})
	}
}
