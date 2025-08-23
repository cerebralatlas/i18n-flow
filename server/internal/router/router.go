package router

import (
	v1 "i18n-flow/internal/api/v1"
	"i18n-flow/internal/middleware"
	"i18n-flow/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		// ❌ 公开接口：不需要 JWT
		apiV1.POST("/register", v1.Register)
		apiV1.POST("/login", v1.Login)

		// ✅ 受保护接口：需要 JWT 验证
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetString("userID")
				response.Success(c, gin.H{"user_id": userID})
			})
		}
	}

	return r
}
