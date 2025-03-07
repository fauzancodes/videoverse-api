package service

import (
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

func CreatePlaylist(userID string, request dto.PlaylistRequest) (response models.VAPlaylist, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VAPlaylist{
		Title:       request.Title,
		Description: request.Description,
		UserID:      &parsedUserUUID,
		Status:      request.Status,
		Visibility:  request.Visibility,
	}

	response, err = repository.CreatePlaylist(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode, err = AddVideosToPlaylist(request.VideoIDs, &response)
	if err != nil {
		return
	}

	statusCode = http.StatusCreated
	return
}

func AddVideosToPlaylist(videoIDs []string, playlist *models.VAPlaylist) (statusCode int, err error) {
	var parsedVideoIDs []uuid.UUID
	for _, item := range videoIDs {
		var parsedVideoUUID uuid.UUID
		parsedVideoUUID, err = uuid.Parse(item)
		if err != nil {
			err = fmt.Errorf("failed to parse video UUID %s: %s", item, err.Error())
			statusCode = http.StatusBadRequest
			return
		}

		parsedVideoIDs = append(parsedVideoIDs, parsedVideoUUID)
	}

	videos, err := repository.GetVideoByIDs(parsedVideoIDs, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get videos: %s", err.Error())

		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.AddVideosToPlaylist(videos, *playlist)
	if err != nil {
		err = fmt.Errorf("failed to add videos into playlist: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	playlist.Videos = videos

	return
}

func GetPlaylistByID(id string, preloadFields []string) (data models.VAPlaylist, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err = repository.GetPlaylistByID(parsedUUID, preloadFields)
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

func GetPlaylists(visibility, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAPlaylist, statusCode int, err error) {
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
	if param.Custom != "" {
		filter += " AND status = ?"
		filterValues = append(filterValues, param.Custom.(string))
	}
	if param.Search != "" {
		filter += " AND (title ILIKE ? OR description ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetPlaylists(dto.FindParameter{
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

func UpdatePlaylist(id string, request dto.PlaylistRequest) (response models.VAPlaylist, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetPlaylistByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
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
	if request.Visibility != "" {
		data.Visibility = request.Visibility
	}
	data.Status = request.Status

	response, err = repository.UpdatePlaylist(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if len(request.VideoIDs) > 0 {
		err = repository.ClearVideosFromPlaylist(data)
		if err != nil {
			err = fmt.Errorf("failed to clear videos in playlist: %s", err.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		statusCode, err = AddVideosToPlaylist(request.VideoIDs, &response)
		if err != nil {
			return
		}
	}

	statusCode = http.StatusOK
	return
}

func DeletePlaylist(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetPlaylistByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.ClearVideosFromPlaylist(data)
	if err != nil {
		err = fmt.Errorf("failed to clear videos in playlist: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeletePlaylist(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
