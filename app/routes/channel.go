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

		subscriptions := channels.Group("/subscriptions")
		{
			subscriptions.POST("", middlewares.Auth(), controllers.CreateSubscription)
			subscriptions.GET("", controllers.GetSubscriptions)
			subscriptions.DELETE("/:channel_id", middlewares.Auth(), controllers.DeleteSubscription)
		}
	}
}
