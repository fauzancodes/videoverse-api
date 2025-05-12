package controllers

import (
	"fmt"
	"net/http"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/upload"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	buf, statusCode, err := upload.GetRemoteFile("/assets/html/index.html")
	if err != nil {
		c.JSON(
			statusCode, dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
		return
	}

	c.Data(http.StatusOK, "text/html", buf.Bytes())
}

func DownloadPostmanCollection(c *gin.Context) {
	fmt.Println("Triggered!")
	buf, statusCode, err := upload.GetRemoteFile("/docs/VideoVerse API.postman_collection.json")
	if err != nil {
		fmt.Println("err:", err.Error())
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="VideoVerse API.postman_collection.json"`)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func DownloadPostmanEnvironment(c *gin.Context) {
	buf, statusCode, err := upload.GetRemoteFile("/docs/VideoVerse API.postman_environment.json")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="VideoVerse API.postman_environment.json"`)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}
