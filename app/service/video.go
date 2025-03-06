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
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	parsedCategoryUUID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		err = errors.New("failed to parse category UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	tags, err := json.Marshal(request.Tags)
	if err != nil {
		err = errors.New("failed to marshal tags: " + err.Error())
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
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoByID(id string, preloadFields []string) (data models.VAVideo, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err = repository.GetVideoByID(parsedUUID, preloadFields)
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
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
		err = errors.New("failed to get data: " + err.Error())
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

func UpdateVideo(id string, request dto.VideoRequest) (response models.VAVideo, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetVideoByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
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
			err = errors.New("failed to parse category UUID: " + err.Error())
			statusCode = http.StatusInternalServerError
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
			err = errors.New("failed to marshal tags: " + err.Error())
			statusCode = http.StatusBadRequest
			return
		}
		data.Tags = string(tags)
	}
	data.Status = request.Status

	response, err = repository.UpdateVideo(data)
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteVideo(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetVideoByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteVideo(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func CreateVideoCategory(userID string, request dto.VideoCategoryRequest) (response models.VAVideoCategory, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
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
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoCategoryByID(id string, preloadFields []string) (data models.VAVideoCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err = repository.GetVideoCategoryByID(parsedUUID, preloadFields)
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
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
		err = errors.New("failed to get data: " + err.Error())
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

func UpdateVideoCategory(id string, request dto.VideoCategoryRequest) (response models.VAVideoCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetVideoCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
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
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteVideoCategory(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetVideoCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteVideoCategory(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
