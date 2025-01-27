package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func VideoRoute(api *gin.RouterGroup) {
	api.POST("/upload", middlewares.CheckAPIKey(), middlewares.Auth(), controllers.Upload)
}
