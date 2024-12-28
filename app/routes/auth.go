package routes

import (
	"github.com/fauzancodes/videoverse-api/app/controllers"
	"github.com/fauzancodes/videoverse-api/app/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoute(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", middlewares.CheckAPIKey(), middlewares.Auth(), controllers.GetCurrentUser)
		auth.PATCH("/update-account", middlewares.CheckAPIKey(), middlewares.Auth(), controllers.UpdateAccount)
		auth.DELETE("/delete-account", middlewares.CheckAPIKey(), middlewares.Auth(), controllers.DeleteAccount)

		emailVerfication := auth.Group("/email-verification")
		{
			emailVerfication.GET("/:token", controllers.VerifyUser)
			emailVerfication.POST("/resend", middlewares.CheckAPIKey(), controllers.ResendEmailVerification)
		}

		resetPassword := auth.Group("/reset-password")
		{
			resetPassword.POST("/send", middlewares.CheckAPIKey(), controllers.SendForgotPasswordRequest)
			resetPassword.GET("/instruction/:token", controllers.SendResetPasswordRequestInstruction)
			resetPassword.POST("", middlewares.CheckAPIKey(), controllers.ResetPassword)
		}
	}
}
