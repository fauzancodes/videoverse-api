package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func CreateChannel(c *gin.Context) {
	var request dto.ChannelRequest
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

	result, statusCode, err := service.CreateChannel(userID, request)
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
			Data:    result,
		},
	)
}

func GetChannels(c *gin.Context) {
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
	withSubscribers, _ := strconv.ParseBool(c.Query("with_subscribers"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User", "User.Profile")
	}
	if withSubscribers {
		preloadFields = append(preloadFields, "Subscribers", "Subscribers.Subscriber", "Subscribers.Subscriber.Profile")
	}

	fmt.Println("preloadFields: ", preloadFields)

	param := utils.PopulatePaging(c, "")
	data, _, statusCode, err := service.GetChannels(userID, param, preloadFields)
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

func GetChannelByID(c *gin.Context) {
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
	withSubscribers, _ := strconv.ParseBool(c.Query("with_subscribers"))

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User", "User.Profile")
	}
	if withSubscribers {
		preloadFields = append(preloadFields, "Subscribers", "Subscribers.Subscriber", "Subscribers.Subscriber.Profile")
	}

	data, statusCode, err := service.GetChannelByID(id, preloadFields)
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
		http.StatusOK,
		dto.Response{
			Status:  http.StatusOK,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateChannel(c *gin.Context) {
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

	var request dto.ChannelRequest
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

	data, statusCode, err := service.UpdateChannel(id, request)
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

func DeleteChannel(c *gin.Context) {
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

	statusCode, err = service.DeleteChannel(id)
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

func CreateSubscription(c *gin.Context) {
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

	var request dto.SubscriptionRequest
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

	result, statusCode, err := service.CreateSubscription(userID, request)
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

func GetSubscriptions(c *gin.Context) {
	userID, statusCode, err := utils.QueryParamUUID(c, "subscriber_id")
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
	channelID, statusCode, err := utils.QueryParamUUID(c, "channel_id")
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

	withSubscriber, _ := strconv.ParseBool(c.Query("with_subscriber"))
	withChannel, _ := strconv.ParseBool(c.Query("with_channel"))

	var preloadFields []string
	if withSubscriber {
		preloadFields = append(preloadFields, "Subscriber", "Subscriber.Profile")
	}
	if withChannel {
		preloadFields = append(preloadFields, "Channel")
	}

	param := utils.PopulatePaging(c, "")
	data, _, statusCode, err := service.GetSubscriptions(channelID, userID, param, preloadFields)
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

func DeleteSubscription(c *gin.Context) {
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

	channelID, statusCode, err := utils.ParamUUID(c, "channel_id")
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

	statusCode, err = service.DeleteSubscription(channelID, userID)
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
