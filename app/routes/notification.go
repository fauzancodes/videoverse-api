package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func NotificationRoute(api *gin.RouterGroup) {
	notifications := api.Group("/notifications", middlewares.Auth(), middlewares.CheckAPIKey())
	{
		notifications.GET("", controllers.GetNotifications)
		notifications.PATCH("/:id/read", controllers.ReadNotification)
		notifications.PATCH("/:id/unread", controllers.UnreadNotification)
		notifications.DELETE("/:id", controllers.DeleteNotification)
	}
}
