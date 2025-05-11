package service

import (
	"encoding/json"
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

func CreateVideo(userID string, request dto.VideoRequest) (response models.VAVideo, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedCategoryUUID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		err = fmt.Errorf("failed to parse category UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	tags, err := json.Marshal(request.Tags)
	if err != nil {
		err = fmt.Errorf("failed to marshal tags: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VAVideo{
		Title:        request.Title,
		Description:  request.Description,
		CategoryID:   &parsedCategoryUUID,
		VideoUrl:     request.VideoUrl,
		ThumbnailUrl: request.ThumbnailUrl,
		Visibility:   request.Visibility,
		Tags:         string(tags),
		UserID:       &parsedUserUUID,
		Status:       request.Status,
	}

	response, err = repository.CreateVideo(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if request.NotificationRedirect != "" {
		channels, _, _, _ := repository.GetChannels(dto.FindParameter{
			Filter:       "deleted_at IS NULL AND user_id = ?",
			FilterValues: []any{parsedUserUUID},
		}, []string{})

		if len(channels) > 0 {
			subscribers, _, _, _ := repository.GetSubscribtions(dto.FindParameter{
				Filter:       "deleted_at IS NULL AND channel_id = ?",
				FilterValues: []any{channels[0].ID},
			}, []string{})

			for _, subscriber := range subscribers {
				notificationData := dto.NotificationRequest{
					Redirect: request.NotificationRedirect,
				}

				notificationData.Content = fmt.Sprintf("%s uploads new video: %s", channels[0].Name, response.Title)

				_, statusCode, err = CreateNotification(subscriber.SubscriberID.String(), notificationData)
				if err != nil {
					err = fmt.Errorf("failed to create notification: %s", err.Error())
					return
				}
			}
		}
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoByID(id string, preloadFields []string) (data models.VAVideo, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err = repository.GetVideoByID(parsedUUID, preloadFields)
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

func GetVideos(visibility, categoryID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAVideo, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	var baseFilterValues []any
	if userID != "" {
		baseFilter += " AND user_id = ?"
		baseFilterValues = append(baseFilterValues, userID)
	}
	filter := baseFilter
	filterValues := baseFilterValues

	if visibility != "" {
		filter += " AND visibility = ?"
		filterValues = append(filterValues, visibility)
	}
	if categoryID != "" {
		filter += " AND category_id = ?"
		filterValues = append(filterValues, categoryID)
	}
	if param.Custom != "" {
		filter += " AND status = ?"
		filterValues = append(filterValues, param.Custom.(string))
	}
	if param.Search != "" {
		filter += " AND (title ILIKE ? OR description ILIKE ? OR tags ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetVideos(dto.FindParameter{
		BaseFilter:       baseFilter,
		BaseFilterValues: baseFilterValues,
		Filter:           filter,
		FilterValues:     filterValues,
		Limit:            param.Limit,
		Order:            param.Order,
		Offset:           param.Offset,
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

func UpdateVideo(id, userID string, request dto.VideoRequest) (response models.VAVideo, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetVideoByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if *data.UserID != parsedUserUUID {
		err = errors.New("you are not authorized to update this data")
		statusCode = http.StatusForbidden
		return
	}

	if request.Title != "" {
		data.Title = request.Title
	}
	if request.Description != "" {
		data.Description = request.Description
	}
	if request.CategoryID != "" {
		var parsedCategoryUUID uuid.UUID
		parsedCategoryUUID, err = uuid.Parse(request.CategoryID)
		if err != nil {
			err = fmt.Errorf("failed to parse category UUID: %s", err.Error())
			statusCode = http.StatusBadRequest
			return
		}
		data.CategoryID = &parsedCategoryUUID
	}
	if request.VideoUrl != "" {
		data.VideoUrl = request.VideoUrl
	}
	if request.ThumbnailUrl != "" {
		data.ThumbnailUrl = request.ThumbnailUrl
	}
	if request.Visibility != "" {
		data.Visibility = request.Visibility
	}
	if len(request.Tags) > 0 {
		var tags []byte
		tags, err = json.Marshal(request.Tags)
		if err != nil {
			err = fmt.Errorf("failed to marshal tags: %s", err.Error())
			statusCode = http.StatusBadRequest
			return
		}
		data.Tags = string(tags)
	}
	data.Status = request.Status

	response, err = repository.UpdateVideo(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteVideo(id, userID string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetVideoByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if *data.UserID != parsedUserUUID {
		err = errors.New("you are not authorized to delete this data")
		statusCode = http.StatusForbidden
		return
	}

	err = repository.DeleteVideo(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func CreateVideoCategory(userID string, request dto.VideoCategoryRequest) (response models.VAVideoCategory, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VAVideoCategory{
		Title:       request.Title,
		Description: request.Description,
		UserID:      &parsedUserUUID,
		Status:      request.Status,
	}

	response, err = repository.CreateVideoCategory(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoCategoryByID(id string, preloadFields []string) (data models.VAVideoCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err = repository.GetVideoCategoryByID(parsedUUID, preloadFields)
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

func GetVideoCategories(userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAVideoCategory, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	var baseFilterValues []any
	if userID != "" {
		baseFilter += " AND user_id = ?"
		baseFilterValues = append(baseFilterValues, userID)
	}
	filter := baseFilter
	filterValues := baseFilterValues

	if param.Custom != "" {
		filter += " AND status = ?"
		filterValues = append(filterValues, param.Custom.(string))
	}
	if param.Search != "" {
		filter += " AND (title ILIKE ? OR description ILIKE ? OR tags ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetVideoCategories(dto.FindParameter{
		BaseFilter:       baseFilter,
		BaseFilterValues: baseFilterValues,
		Filter:           filter,
		FilterValues:     filterValues,
		Limit:            param.Limit,
		Order:            param.Order,
		Offset:           param.Offset,
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

func UpdateVideoCategory(id, userID string, request dto.VideoCategoryRequest) (response models.VAVideoCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetVideoCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if *data.UserID != parsedUserUUID {
		err = errors.New("you are not authorized to update this data")
		statusCode = http.StatusForbidden
		return
	}

	if request.Title != "" {
		data.Title = request.Title
	}
	if request.Description != "" {
		data.Description = request.Description
	}
	data.Status = request.Status

	response, err = repository.UpdateVideoCategory(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteVideoCategory(id, userID string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetVideoCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if *data.UserID != parsedUserUUID {
		err = errors.New("you are not authorized to delete this data")
		statusCode = http.StatusForbidden
		return
	}

	err = repository.DeleteVideoCategory(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
