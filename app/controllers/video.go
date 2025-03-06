package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func CreateVideo(c *gin.Context) {
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

	var request dto.VideoRequest
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

	result, statusCode, err := service.CreateVideo(userID, request)
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

func GetVideos(c *gin.Context) {
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

	visibility := c.Query("visibility")
	categoryID, statusCode, err := utils.QueryParamUUID(c, "category_id")
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

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetVideos(visibility, categoryID, userID, param, preloadFields)
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

func GetPublicVideos(c *gin.Context) {
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

	visibility := "public"
	categoryID, statusCode, err := utils.QueryParamUUID(c, "category_id")
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

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetVideos(visibility, categoryID, userID, param, preloadFields)
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

func GetVideoByID(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}

	data, statusCode, err := service.GetVideoByID(id, preloadFields)
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

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func GetPublicVideoByID(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}

	data, statusCode, err := service.GetVideoByID(id, preloadFields)
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

	if strings.EqualFold(data.Visibility, "private") {
		c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  http.StatusNotFound,
				Message: "Failed to get data",
				Error:   "failed to get data",
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateVideo(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	var request dto.VideoRequest
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

	data, statusCode, err := service.UpdateVideo(id, request)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to update data",
			Data:    data,
		},
	)
}

func DeleteVideo(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	statusCode, err = service.DeleteVideo(id)
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

func CreateVideoCategory(c *gin.Context) {
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

	var request dto.VideoCategoryRequest
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

	result, statusCode, err := service.CreateVideoCategory(userID, request)
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

func GetVideoCategories(c *gin.Context) {
	userID, statusCode, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get current user ID",
				Error:   err.Error(),
			},
		)
		return
	}

	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetVideoCategories(userID, param, preloadFields)
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

func GetVideoCategoryByID(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	data, statusCode, err := service.GetVideoCategoryByID(id, preloadFields)
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

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateVideoCategory(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	var request dto.VideoCategoryRequest
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

	data, statusCode, err := service.UpdateVideoCategory(id, request)
	if err != nil {
		c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to update data",
			Data:    data,
		},
	)
}

func DeleteVideoCategory(c *gin.Context) {
	id, statusCode, err := utils.ParamUUID(c, "id")
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

	statusCode, err = service.DeleteVideoCategory(id)
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
