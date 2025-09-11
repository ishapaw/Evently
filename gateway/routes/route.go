package routes

import (
	"gateway/auth"
	"gateway/proxy"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.Any("/users/*path", proxy.ReverseProxy("http://users-service:8081"))

	protected := api.Group("/v1")
	protected.Use(auth.AuthMiddleware())
	{
		protected.Any("/events/*path", proxy.ReverseProxy("http://events-service:8082"))
	}
}
