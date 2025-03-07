package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/bcrypt"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(request dto.UserRequest) (response models.VAUser, statusCode int, err error) {
	data := models.VAUser{
		Email:    request.Email,
		Password: bcrypt.HashPassword(request.Password),
	}

	response, err = repository.CreateUser(data)
	if err != nil {
		err = fmt.Errorf("failed to create data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetUserByID(id string, preloadFields []string) (data models.VAUser, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	data, err = repository.GetUserByID(parsedUUID, preloadFields)
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

func GetUsers(email string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAUser, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if email != "" {
		filter += " AND email = ?"
		filterValues = append(filterValues, email)
	}
	if param.Search != "" {
		filter += " AND email ILIKE ?"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetUsers(dto.FindParameter{
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

func UpdateUser(id string, request dto.UserRequest) (response models.VAUser, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}
	data, err := repository.GetUserByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if request.Email != "" {
		data.Email = request.Email
	}
	if request.Password != "" {
		data.Password = bcrypt.HashPassword(request.Password)
	}

	response, err = repository.UpdateUser(data)
	if err != nil {
		err = fmt.Errorf("failed to update data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteUser(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID: %s", err.Error())
		statusCode = http.StatusBadRequest
		return
	}

	data, err := repository.GetUserByID(parsedUUID, []string{})
	if err != nil {
		err = fmt.Errorf("failed to get data: %s", err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteUser(data)
	if err != nil {
		err = fmt.Errorf("failed to delete data: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
