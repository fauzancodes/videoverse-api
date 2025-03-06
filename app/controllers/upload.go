package controllers

import (
	"net/http"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	userID, statusCode, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get current userID",
				Error:   err.Error(),
			},
		)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  500,
				Message: "Failed to get file from form",
				Error:   err.Error(),
			},
		)
		return
	}

	responseURL, statusCode, err := service.Upload(file, userID)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to upload product picture",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to upload",
			Data:    responseURL,
		},
	)
}
