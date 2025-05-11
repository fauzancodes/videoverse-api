package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateChannel(userID string, request dto.ChannelRequest) (response models.VAChannel, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	checkData, _, _, _ := repository.GetChannels(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND user_id = ?",
		FilterValues: []any{parsedUserUUID},
	}, []string{})
	if len(checkData) > 0 {
		err = errors.New("channel for this user has been created")
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VAChannel{
		Name:        request.Name,
		Picture:     request.Picture,
		Description: request.Description,
		Location:    request.Location,
		UserID:      parsedUserUUID,
	}

	response, err = repository.CreateChannel(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetChannelByID(id string, preloadFields []string) (data models.VAChannel, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err = repository.GetChannelByID(parsedUUID, preloadFields)
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func GetChannels(userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAChannel, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if userID != "" {
		filter += " AND user_id = ?"
		filterValues = append(filterValues, userID)
	}
	if param.Search != "" {
		filter += " AND (name ILIKE ? OR description ILIKE ? OR location ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetChannels(dto.FindParameter{
		BaseFilter:   baseFilter,
		Filter:       filter,
		FilterValues: filterValues,
		Limit:        param.Limit,
		Order:        param.Order,
		Offset:       param.Offset,
	}, preloadFields)
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	statusCode = http.StatusOK
	return
}

func UpdateChannel(id string, request dto.ChannelRequest) (response models.VAChannel, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	data, err := repository.GetChannelByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if request.Name != "" {
		data.Name = request.Name
	}
	if request.Description != "" {
		data.Description = request.Description
	}
	if request.Location != "" {
		data.Location = request.Location
	}
	if request.Picture != "" {
		data.Picture = request.Picture
	}

	response, err = repository.UpdateChannel(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteChannel(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetChannelByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteChannel(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func CreateSubscribtion(userID string, request dto.SubscribtionRequest) (response models.VASubscribtion, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedChannelUUID, err := uuid.Parse(request.ChannelID)
	if err != nil {
		err = fmt.Errorf("failed to parse channel UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	channel, err := repository.GetChannelByID(parsedChannelUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get channel: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	checkData, _, _, _ := repository.GetSubscribtions(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND channel_id = ? AND subscriber_id = ?",
		FilterValues: []any{parsedChannelUUID, parsedUserUUID},
	}, []string{})
	if len(checkData) > 0 {
		err = errors.New("this channel has been subscribed")
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VASubscribtion{
		ChannelID:    parsedChannelUUID,
		SubscriberID: parsedUserUUID,
	}

	response, err = repository.CreateSubscribtion(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if request.NotificationRedirect != "" {
		var user models.VAUser
		user, err = repository.GetUserByID(parsedUserUUID, []string{"Profile"})
		if err != nil {
			err = fmt.Errorf("failed to get user for notification: %s", err.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		notificationData := dto.NotificationRequest{
			Redirect: request.NotificationRedirect,
		}
		userIdentifier := user.Email
		if user.Profile.Firstname != "" {
			userIdentifier = user.Profile.Firstname
		}
		if user.Profile.Firstname != "" && user.Profile.Lastname != "" {
			userIdentifier = fmt.Sprintf("%s %s", user.Profile.Firstname, user.Profile.Lastname)
		}

		notificationData.Content = fmt.Sprintf("%s subscribes to your channel!", userIdentifier)

		_, statusCode, err = CreateNotification(channel.UserID.String(), notificationData)
		if err != nil {
			err = fmt.Errorf("failed to create notification: %s", err.Error())
			return
		}
	}

	statusCode = http.StatusCreated
	return
}

func GetSubscribtions(channelID, subscriberID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VASubscribtion, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if subscriberID != "" {
		filter += " AND subscriber_id = ?"
		filterValues = append(filterValues, subscriberID)
	}
	if channelID != "" {
		filter += " AND channel_id = ?"
		filterValues = append(filterValues, channelID)
	}

	data, total, totalFiltered, err := repository.GetSubscribtions(dto.FindParameter{
		BaseFilter:   baseFilter,
		Filter:       filter,
		FilterValues: filterValues,
		Limit:        param.Limit,
		Order:        param.Order,
		Offset:       param.Offset,
	}, preloadFields)
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	statusCode = http.StatusOK
	return
}

func DeleteSubscribtion(channelID, userID string, request dto.UnsubscribeRequest) (statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedChannelUUID, err := uuid.Parse(channelID)
	if err != nil {
		err = fmt.Errorf("failed to parse channel UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, _, _, err := repository.GetSubscribtions(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND channel_id = ? AND subscriber_id = ?",
		FilterValues: []any{parsedChannelUUID, parsedUserUUID},
	}, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}
	if len(data) == 0 {
		err = errors.New("failed to get data: data not found")
		statusCode = http.StatusNotFound
		return
	}

	if data[0].SubscriberID != parsedUserUUID {
		err = errors.New("you are not authorized to delete this data")
		statusCode = http.StatusForbidden
		return
	}

	err = repository.DeleteSubscribtion(data[0])
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if request.NotificationRedirect != "" {
		var channel models.VAChannel
		channel, err = repository.GetChannelByID(parsedChannelUUID, []string{})
		if err != nil {
			err = fmt.Errorf("failed to get channel for notification: %s", err.Error())
			statusCode = http.StatusBadRequest
			return
		}

		var user models.VAUser
		user, err = repository.GetUserByID(parsedUserUUID, []string{"Profile"})
		if err != nil {
			err = fmt.Errorf("failed to get user for notification: %s", err.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		notificationData := dto.NotificationRequest{
			Redirect: request.NotificationRedirect,
		}
		userIdentifier := user.Email
		if user.Profile.Firstname != "" {
			userIdentifier = user.Profile.Firstname
		}
		if user.Profile.Firstname != "" && user.Profile.Lastname != "" {
			userIdentifier = fmt.Sprintf("%s %s", user.Profile.Firstname, user.Profile.Lastname)
		}

		notificationData.Content = fmt.Sprintf("%s unsubscribes to your channel!", userIdentifier)

		_, statusCode, err = CreateNotification(channel.UserID.String(), notificationData)
		if err != nil {
			err = fmt.Errorf("failed to create notification: %s", err.Error())
			return
		}
	}

	statusCode = http.StatusOK
	return
}
