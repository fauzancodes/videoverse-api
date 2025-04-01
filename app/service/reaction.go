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

func CreateVideoLike(userID string, request dto.VideoLikeRequest) (response models.VAVideoLike, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedVideoUUID, err := uuid.Parse(request.VideoID)
	if err != nil {
		err = fmt.Errorf("failed to parse video UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	_, err = repository.GetVideoByID(parsedVideoUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get video: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	checkData, _, _, _ := repository.GetVideoLikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
	}, []string{})
	if len(checkData) > 0 {
		err = errors.New("this video has been liked")
		statusCode = http.StatusBadRequest
		return
	}

	dislikeData, _, _, err := repository.GetVideoDislikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
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
	if len(dislikeData) != 0 {
		repository.DeleteVideoDislike(dislikeData[0])
	}

	data := models.VAVideoLike{
		VideoID: parsedVideoUUID,
		UserID:  parsedUserUUID,
	}

	response, err = repository.CreateVideoLike(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoLikes(videoID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAVideoLike, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if userID != "" {
		filter += " AND user_id = ?"
		filterValues = append(filterValues, userID)
	}
	if videoID != "" {
		filter += " AND video_id = ?"
		filterValues = append(filterValues, videoID)
	}

	data, total, totalFiltered, err := repository.GetVideoLikes(dto.FindParameter{
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

func DeleteVideoLike(videoID, userID string) (statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedVideoUUID, err := uuid.Parse(videoID)
	if err != nil {
		err = fmt.Errorf("failed to parse video UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, _, _, err := repository.GetVideoLikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
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

	if data[0].UserID != parsedUserUUID {
		err = errors.New("you are not authorized to delete this data")
		statusCode = http.StatusForbidden
		return
	}

	err = repository.DeleteVideoLike(data[0])
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func CreateVideoDislike(userID string, request dto.VideoDislikeRequest) (response models.VAVideoDislike, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedVideoUUID, err := uuid.Parse(request.VideoID)
	if err != nil {
		err = fmt.Errorf("failed to parse video UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	_, err = repository.GetVideoByID(parsedVideoUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get video: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	checkData, _, _, _ := repository.GetVideoDislikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
	}, []string{})
	if len(checkData) > 0 {
		err = errors.New("this video has been disliked")
		statusCode = http.StatusBadRequest
		return
	}

	likeData, _, _, err := repository.GetVideoLikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
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
	if len(likeData) != 0 {
		repository.DeleteVideoLike(likeData[0])
	}

	data := models.VAVideoDislike{
		VideoID: parsedVideoUUID,
		UserID:  parsedUserUUID,
	}

	response, err = repository.CreateVideoDislike(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetVideoDislikes(videoID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAVideoDislike, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if userID != "" {
		filter += " AND user_id = ?"
		filterValues = append(filterValues, userID)
	}
	if videoID != "" {
		filter += " AND video_id = ?"
		filterValues = append(filterValues, videoID)
	}

	data, total, totalFiltered, err := repository.GetVideoDislikes(dto.FindParameter{
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

func DeleteVideoDislike(videoID, userID string) (statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	parsedVideoUUID, err := uuid.Parse(videoID)
	if err != nil {
		err = fmt.Errorf("failed to parse video UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, _, _, err := repository.GetVideoDislikes(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND video_id = ? AND user_id = ?",
		FilterValues: []any{parsedVideoUUID, parsedUserUUID},
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

	if data[0].UserID != parsedUserUUID {
		err = errors.New("you are not authorized to delete this data")
		statusCode = http.StatusForbidden
		return
	}

	err = repository.DeleteVideoDislike(data[0])
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
