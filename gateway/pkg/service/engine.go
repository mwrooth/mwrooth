package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/mwrooth/gateway/openapi/go"
	"github.com/mwrooth/gateway/pkg/middleware"
)

const apiPrefix = "/api"

func createOpenAPIHandler(_ string, handler openapi.ContextHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		handler(c.Request.Context(), c)
	}
}

func (s *Service) NewEngine() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()

	apiGroup := engine.Group(apiPrefix)

	apiGroup.Use(middleware.CORSMiddleware())

	for _, route := range openapi.CreateRoutes(s) {
		openAPIHandler := createOpenAPIHandler(route.Pattern, route.HandlerFunc)
		switch route.Method {
		case http.MethodGet:
			apiGroup.GET(route.Pattern, openAPIHandler)
		case http.MethodPost:
			apiGroup.POST(route.Pattern, openAPIHandler)
		case http.MethodPut:
			apiGroup.PUT(route.Pattern, openAPIHandler)
		case http.MethodPatch:
			apiGroup.PATCH(route.Pattern, openAPIHandler)
		case http.MethodDelete:
			apiGroup.DELETE(route.Pattern, openAPIHandler)
		}
	}

	return engine
}
