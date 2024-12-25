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
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusBadRequest,
				Message: "No jwt token provided",
			})

			return
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwt.DecodeToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusUnauthorized,
				Message: "Failed to decode jwt token",
			})

			return
		}

		c.Set("currentUser", claims)
		c.Next()
	}
}
