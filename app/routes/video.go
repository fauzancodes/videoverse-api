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

		playlist := videos.Group("/playlists")
		{
			playlist.POST("", middlewares.Auth(), controllers.CreatePlaylist)
			playlist.GET("", middlewares.Auth(), controllers.GetPlaylists)
			playlist.GET("/public", controllers.GetPublicPlaylists)
			playlist.GET("/public/:id", controllers.GetPublicPlaylistByID)
			playlist.GET("/:id", middlewares.Auth(), controllers.GetPlaylistByID)
			playlist.PATCH("/:id", middlewares.Auth(), controllers.UpdatePlaylist)
			playlist.DELETE("/:id", middlewares.Auth(), controllers.DeletePlaylist)
		}

		likes := videos.Group("/likes")
		{
			likes.POST("", middlewares.Auth(), controllers.CreateVideoLike)
			likes.GET("", controllers.GetVideoLikes)
			likes.DELETE("/:video_id", middlewares.Auth(), controllers.DeleteVideoLike)
		}

		dislikes := videos.Group("/dislikes")
		{
			dislikes.POST("", middlewares.Auth(), controllers.CreateVideoDislike)
			dislikes.GET("", controllers.GetVideoDislikes)
			dislikes.DELETE("/:video_id", middlewares.Auth(), controllers.DeleteVideoDislike)
		}

		comment := videos.Group("/comments")
		{
			comment.POST("", middlewares.Auth(), controllers.CreateComment)
			comment.GET("", controllers.GetComments)
			comment.GET("/:id", controllers.GetCommentByID)
			comment.PATCH("/:id", middlewares.Auth(), controllers.UpdateComment)
			comment.DELETE("/:id", middlewares.Auth(), controllers.DeleteComment)
		}
	}
}
