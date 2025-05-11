package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/veezify/interview-ai-backend/docs"
	"github.com/veezify/interview-ai-backend/internal/api/handler"
	"github.com/veezify/interview-ai-backend/internal/api/middleware"
	"github.com/veezify/interview-ai-backend/internal/service"
)

func SetupRouter(
	authService *service.AuthService,
) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.CORS())

	authHandler := handler.NewAuthHandler(authService)

	public := r.Group("/api")
	{
		public.POST("/auth/login", authHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authService))
	{

		protected.GET("/me", authHandler.GetMe)

		//associations := protected.Group("/associations")
		//{
		//
		//	associations.GET("")
		//
		//	association := associations.Group("/:associationID")
		//	{
		//
		//		association.GET("", middleware.HasPermission("view_association"))
		//		association.PUT("", middleware.HasPermission("edit_association"))
		//	}
		//}
	}

	return r
}
