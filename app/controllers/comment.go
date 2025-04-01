package controllers

import (
	"net/http"
	"strconv"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
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

	var request dto.CommentRequest
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

	result, statusCode, err := service.CreateComment(userID, request)
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

func GetComments(c *gin.Context) {
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
	parentID, statusCode, err := utils.QueryParamUUID(c, "parent_id")
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
	withVideo, _ := strconv.ParseBool(c.Query("with_video"))
	withParent, _ := strconv.ParseBool(c.Query("with_parent"))
	withReplies, _ := strconv.ParseBool(c.Query("with_replies"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideo {
		preloadFields = append(preloadFields, "Video")
	}
	if withParent {
		preloadFields = append(preloadFields, "Parent")
	}
	if withReplies {
		preloadFields = append(preloadFields, "Replies", "Replies.User")
	}

	param := utils.PopulatePaging(c, "")
	data, _, statusCode, err := service.GetComments(parentID, videoID, userID, param, preloadFields)
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

func GetCommentByID(c *gin.Context) {
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
	withVideo, _ := strconv.ParseBool(c.Query("with_video"))
	withParent, _ := strconv.ParseBool(c.Query("with_parent"))
	withReplies, _ := strconv.ParseBool(c.Query("with_replies"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withVideo {
		preloadFields = append(preloadFields, "Video")
	}
	if withParent {
		preloadFields = append(preloadFields, "Parent")
	}
	if withReplies {
		preloadFields = append(preloadFields, "Replies", "Replies.User")
	}

	data, statusCode, err := service.GetCommentByID(id, preloadFields)
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

func UpdateComment(c *gin.Context) {
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

	var request dto.CommentUpdateRequest
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

	data, statusCode, err := service.UpdateComment(userID, id, request)
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

func DeleteComment(c *gin.Context) {
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

	statusCode, err = service.DeleteComment(userID, id)
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
