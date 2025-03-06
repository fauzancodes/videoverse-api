package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func VideoRoute(api *gin.RouterGroup) {
	api.POST("/upload", middlewares.CheckAPIKey(), middlewares.Auth(), controllers.Upload)

	videos := api.Group("/videos", middlewares.CheckAPIKey())
	{
		categories := videos.Group("/categories", middlewares.Auth())
		{
			categories.POST("", controllers.CreateVideoCategory)
			categories.GET("", controllers.GetVideoCategories)
			categories.GET("/:id", controllers.GetVideoCategoryByID)
			categories.PATCH("/:id", controllers.UpdateVideoCategory)
			categories.DELETE("/:id", controllers.DeleteVideoCategory)
		}

		videos.POST("", middlewares.Auth(), controllers.CreateVideo)
		videos.GET("", middlewares.Auth(), controllers.GetVideos)
		videos.GET("/public", controllers.GetPublicVideos)
		videos.GET("/public/:id", controllers.GetPublicVideoByID)
		videos.GET("/:id", middlewares.Auth(), controllers.GetVideoByID)
		videos.PATCH("/:id", middlewares.Auth(), controllers.UpdateVideo)
		videos.DELETE("/:id", middlewares.Auth(), controllers.DeleteVideo)
	}
}
