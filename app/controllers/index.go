package controllers

import (
	"net/http"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/upload"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	buf, statusCode, err := upload.GetRemoteFile("/assets/html/index.html")
	if err != nil {
		c.AbortWithStatusJSON(statusCode, dto.Response{
			Status:  statusCode,
			Message: "Failed to get file",
			Error:   err.Error(),
		})
	}

	c.Data(http.StatusOK, "text/html", buf.Bytes())
}

func DownloadPostmanCollection(c *gin.Context) {
	buf, statusCode, err := upload.GetRemoteFile("/docs/Sales Demo API.postman_collection.json")
	if err != nil {
		c.AbortWithStatusJSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
	}

	c.Header("Content-Disposition", `attachment; filename="Sales Demo API.postman_collection.json"`)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func DownloadPostmanEnvironment(c *gin.Context) {
	buf, statusCode, err := upload.GetRemoteFile("/docs/Sales Demo API.postman_environment.json")
	if err != nil {
		c.AbortWithStatusJSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
	}

	c.Header("Content-Disposition", `attachment; filename="Sales Demo API.postman_environment.json"`)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}
