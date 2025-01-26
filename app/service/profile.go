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

func CreateProfile(userID string, request dto.ProfileRequest) (response models.VAProfile, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data := models.VAProfile{
		Firstname:   request.Firstname,
		Lastname:    request.Lastname,
		Gender:      request.Gender,
		Picture:     request.Picture,
		Description: request.Description,
		Location:    request.Location,
		UserID:      parsedUserUUID,
	}

	response, err = repository.CreateProfile(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if len(request.SocialMedia) > 0 {
		var socialMediaInput []models.VASocialMedia
		for _, item := range request.SocialMedia {
			socialMediaInput = append(socialMediaInput, models.VASocialMedia{
				Name:      item.Name,
				Link:      item.Link,
				ProfileID: response.ID,
			})
		}

		response.SocialMedia, err = repository.CreateSocialMediaMulti(socialMediaInput)
		if err != nil {
			err = errors.New("failed to create social media: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}
	}

	statusCode = http.StatusCreated
	return
}

func GetProfileByID(id string, preloadFields []string) (data models.VAProfile, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err = repository.GetProfileByID(parsedUUID, preloadFields)
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

func GetProfiles(userID, gender string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.VAProfile, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if userID != "" {
		filter += " AND user_id = ?"
		filterValues = append(filterValues, userID)
	}
	if gender != "" {
		filter += " AND gender = ?"
		filterValues = append(filterValues, gender)
	}
	if param.Search != "" {
		filter += " AND (firstname ILIKE ? OR lastname ILIKE ? OR description ILIKE ? OR location ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search), fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetProfiles(dto.FindParameter{
		BaseFilter:   baseFilter,
		Filter:       filter,
		FilterValues: filterValues,
		Limit:        param.Limit,
		Order:        param.Order,
		Offset:       param.Offset,
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

func UpdateProfile(id string, request dto.ProfileRequest) (response models.VAProfile, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	data, err := repository.GetProfileByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if request.Firstname != "" {
		data.Firstname = request.Firstname
	}
	if request.Lastname != "" {
		data.Lastname = request.Lastname
	}
	if request.Gender != "" {
		data.Gender = request.Gender
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

	if len(request.SocialMedia) > 0 {
		err = repository.DeleteSocialMediaByProfileID(data.ID)
		if err != nil {
			err = errors.New("failed to delete social media data: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		var socialMediaInput []models.VASocialMedia
		for _, item := range request.SocialMedia {
			socialMediaInput = append(socialMediaInput, models.VASocialMedia{
				Name:      item.Name,
				Link:      item.Link,
				ProfileID: response.ID,
			})
		}

		response.SocialMedia, err = repository.CreateSocialMediaMulti(socialMediaInput)
		if err != nil {
			err = errors.New("failed to create social media: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}
	}

	response, err = repository.UpdateProfile(data)
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteProfile(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProfileByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteSocialMediaByProfileID(data.ID)
	if err != nil {
		err = errors.New("failed to delete social media data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteProfile(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
