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

func CreateComment(userID string, request dto.CommentRequest) (response models.VAComment, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	var parentUUID *uuid.UUID
	if request.ParentID != "" {
		var parsedParentUUID uuid.UUID
		parsedParentUUID, err = uuid.Parse(request.ParentID)
		if err != nil {
			err = fmt.Errorf("failed to parse parent UUID: %s", err.Error())
			statusCode = http.StatusBadRequest
			return
		}
		parentUUID = &parsedParentUUID
	}
	parsedVideoUUID, err := uuid.Parse(request.VideoID)
	if err != nil {
		err = fmt.Errorf("failed to parse video UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data := models.VAComment{
		ParentID: parentUUID,
		UserID:   parsedUserUUID,
		VideoID:  parsedVideoUUID,
		Content:  request.Content,
	}

	response, err = repository.CreateComment(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetCommentByID(id string, preloadFields []string) (data models.VAComment, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err = repository.GetCommentByID(parsedUUID, preloadFields)
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

func GetComments(parentID, videoID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAComment, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if userID != "" {
		filter += " AND user_id = ?"
		filterValues = append(filterValues, userID)
	}
	if parentID != "" {
		filter += " AND parent_id = ?"
		filterValues = append(filterValues, parentID)
	}
	if videoID != "" {
		filter += " AND video_id = ?"
		filterValues = append(filterValues, videoID)
	}
	if param.Search != "" {
		filter += " AND content ILIKE ?"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetComments(dto.FindParameter{
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

func UpdateComment(userID, id string, request dto.CommentUpdateRequest) (response models.VAComment, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetCommentByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	if data.UserID != parsedUserUUID {
		err = fmt.Errorf("you are not allowed to update this comment")
		statusCode = http.StatusBadRequest
		return
	}

	if request.Content != "" {
		data.Content = request.Content
	}

	response, err = repository.UpdateComment(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteComment(userID, id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetCommentByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = fmt.Errorf("failed to parse user UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	if data.UserID != parsedUserUUID {
		err = fmt.Errorf("you are not allowed to delete this comment")
		statusCode = http.StatusBadRequest
		return
	}

	children, _, _, _ := repository.GetComments(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND parent_id = ?",
		FilterValues: []any{data.ID},
	}, []string{})

	if len(children) > 0 {
		for _, child := range children {
			err = repository.DeleteComment(child)
			if err != nil {
				err = fmt.Errorf("failed to delete child: %s", err.Error())
				statusCode = http.StatusInternalServerError
				return
			}
		}
	}

	err = repository.DeleteComment(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
