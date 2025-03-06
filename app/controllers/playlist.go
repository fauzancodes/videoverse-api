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

func CreatePlaylist(c *gin.Context) {
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

	var request dto.PlaylistRequest
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

	result, statusCode, err := service.CreatePlaylist(userID, request)
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

func GetPlaylists(c *gin.Context) {
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

	withUser, _ := strconv.ParseBool(c.Query("with_user"))
	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetPlaylists(visibility, userID, param, preloadFields)
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

func GetPublicPlaylists(c *gin.Context) {
	visibility := "public"
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

	withUser, _ := strconv.ParseBool(c.Query("with_user"))
	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetPlaylists(visibility, userID, param, preloadFields)
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

func GetPlaylistByID(c *gin.Context) {
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
	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	data, statusCode, err := service.GetPlaylistByID(id, preloadFields)
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

func GetPublicPlaylistByID(c *gin.Context) {
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
	withVideos, _ := strconv.ParseBool(c.Query("with_videos"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideos {
		preloadFields = append(preloadFields, "Videos")
	}

	data, statusCode, err := service.GetPlaylistByID(id, preloadFields)
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

func UpdatePlaylist(c *gin.Context) {
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

	var request dto.PlaylistRequest
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

	data, statusCode, err := service.UpdatePlaylist(id, request)
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

func DeletePlaylist(c *gin.Context) {
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

	statusCode, err = service.DeletePlaylist(id)
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
