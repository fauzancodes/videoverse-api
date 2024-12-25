package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoute(app *gin.Engine) {
	auth := app.Group("/auth", middlewares.CheckAPIKey())
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", middlewares.Auth(), controllers.GetCurrentUser)
		auth.PATCH("/update-account", middlewares.Auth(), controllers.UpdateAccount)
		auth.DELETE("/delete-account", middlewares.Auth(), controllers.DeleteAccount)

		emailVerfication := auth.Group("/email-verification")
		{
			emailVerfication.GET("/:token", controllers.VerifyUser)
			emailVerfication.POST("/resend", controllers.ResendEmailVerification)
		}

		resetPassword := auth.Group("/reset-password")
		{
			resetPassword.POST("/send", controllers.SendForgotPasswordRequest)
			resetPassword.GET("/instruction/:token", controllers.SendResetPasswordRequestInstruction)
			resetPassword.POST("", controllers.ResetPassword)
		}
	}
}
