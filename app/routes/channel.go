package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func ChannelRoute(api *gin.RouterGroup) {
	channels := api.Group("/channels", middlewares.CheckAPIKey())
	{
		channels.POST("", middlewares.Auth(), controllers.CreateChannel)
		channels.GET("", controllers.GetChannels)
		channels.GET("/:id", controllers.GetChannelByID)
		channels.PATCH("/:id", middlewares.Auth(), controllers.UpdateChannel)
		channels.DELETE("/:id", middlewares.Auth(), controllers.DeleteChannel)

		subscribtions := channels.Group("/subscribtions")
		{
			subscribtions.POST("", middlewares.Auth(), controllers.CreateSubscribtion)
			subscribtions.GET("", controllers.GetSubscribtions)
			subscribtions.DELETE("/:channel_id", middlewares.Auth(), controllers.DeleteSubscribtion)
		}
	}
}
