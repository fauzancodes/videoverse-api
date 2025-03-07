package controllers

import (
	"net/http"
	"strconv"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func CreateVideoLike(c *gin.Context) {
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

	var request dto.VideoLikeRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
		return
	}

	result, statusCode, err := service.CreateVideoLike(userID, request)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to create",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to create",
			Data:    result,
		},
	)
}

func GetVideoLikes(c *gin.Context) {
	userID, statusCode, err := utils.QueryParamUUID(c, "user_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}
	videoID, statusCode, err := utils.QueryParamUUID(c, "video_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}

	withUser, _ := strconv.ParseBool(c.Query("with_user"))
	withVideo, _ := strconv.ParseBool(c.Query("with_Video"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideo {
		preloadFields = append(preloadFields, "Video")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetVideoLikes(videoID, userID, param, preloadFields)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(statusCode, data)
}

func DeleteVideoLike(c *gin.Context) {
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

	videoID, statusCode, err := utils.ParamUUID(c, "video_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}

	statusCode, err = service.DeleteVideoLike(videoID, userID)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to delete data",
		},
	)
}

func CreateVideoDislike(c *gin.Context) {
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

	var request dto.VideoDislikeRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
		return
	}

	result, statusCode, err := service.CreateVideoDislike(userID, request)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to create",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to create",
			Data:    result,
		},
	)
}

func GetVideoDislikes(c *gin.Context) {
	userID, statusCode, err := utils.QueryParamUUID(c, "user_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}
	videoID, statusCode, err := utils.QueryParamUUID(c, "video_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}

	withUser, _ := strconv.ParseBool(c.Query("with_user"))
	withVideo, _ := strconv.ParseBool(c.Query("with_Video"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideo {
		preloadFields = append(preloadFields, "Video")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetVideoDislikes(videoID, userID, param, preloadFields)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(statusCode, data)
}

func DeleteVideoDislike(c *gin.Context) {
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

	videoID, statusCode, err := utils.ParamUUID(c, "video_id")
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Invalid parameter",
				Error:   err.Error(),
			},
		)
		return
	}

	statusCode, err = service.DeleteVideoDislike(videoID, userID)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to delete data",
		},
	)
}
