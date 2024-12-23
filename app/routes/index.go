package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func Route(app *gin.Engine) {
	app.Static("/assets", "assets")
	app.Static("/docs", "docs")

	app.GET("/", middlewares.StripHTMLMiddleware(), controllers.Index)
	app.GET("/postman/collection", middlewares.StripHTMLMiddleware(), controllers.DownloadPostmanCollection)
	app.GET("/postman/environment", middlewares.StripHTMLMiddleware(), controllers.DownloadPostmanEnvironment)
}
