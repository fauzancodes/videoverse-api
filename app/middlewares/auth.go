package middlewares

import (
	"net/http"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve token from http header
		token := c.GetHeader("Authorization")

		//check if token is in http header
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusBadRequest,
				Message: "No jwt token provided",
			})
		}

		//decode token
		token = strings.Split(token, " ")[1]
		claims, err := jwt.DecodeToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusUnauthorized,
				Message: "Failed to decode jwt token",
				Error:   err.Error(),
			})
		}

		c.Set("currentUser", claims)
		c.Next()
	}
}
