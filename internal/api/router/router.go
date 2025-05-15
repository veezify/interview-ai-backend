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
	interviewService *service.InterviewService,
	userService *service.UserService,
) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.CORS())

	authHandler := handler.NewAuthHandler(authService)
	interviewHandler := handler.NewInterviewHandler(interviewService, userService)
	livekitAuthHandler := handler.NewAuthLivekitHandler()

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

		livekit := protected.Group("/livekit")
		{
			livekit.POST("/token", livekitAuthHandler.GenerateLivekitToken)
		}

		protected.GET("/me", authHandler.GetMe)

		interviews := protected.Group("/interviews")
		{
			interviews.GET("/options", interviewHandler.GetSelectOptions)
			interviews.POST("", interviewHandler.CreateInterview)
			interviews.GET("", interviewHandler.ListInterviews)
			interviews.GET("/:id", interviewHandler.GetInterview)
		}
	}

	return r
}
