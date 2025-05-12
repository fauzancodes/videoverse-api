package handler //for vercel
// package main

import (
	"log"
	"net/http"

	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/fauzancodes/videoverse-api/app/routes"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := Init()

	port := config.LoadConfig().Port

	log.Println("Server: " + config.LoadConfig().BaseUrl + ":" + port)

	app.Run(":" + port)
}

func Main(w http.ResponseWriter, r *http.Request) {
	e := Init()

	e.ServeHTTP(w, r)
}

func Init() *gin.Engine {
	app := gin.Default()

	app.Use(middlewares.Cors())
	app.Use(gzip.Gzip(gzip.BestSpeed))
	app.Use(gin.Logger())
	app.Use(middlewares.Secure())
	app.Use(gin.Recovery())

	config.Database()

	routes.Route(app)

	return app
}
