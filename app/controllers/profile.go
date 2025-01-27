package controllers

import (
	"errors"
	"net/http"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/utils"
	"github.com/fauzancodes/videoverse-api/app/repository"
	"github.com/fauzancodes/videoverse-api/app/service"
	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	var request dto.ProfileRequest
	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
	}

	if len(request.SocialMedia) > 0 {
		for _, item := range request.SocialMedia {
			if err := item.Validate(); err != nil {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					dto.Response{
						Status:  http.StatusBadRequest,
						Message: "Invalid request social media value",
						Error:   err.Error(),
					},
				)
			}
		}
	}

	userID, statusCode, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get current userID",
				Error:   err.Error(),
			},
		)
	}

	var isCreate bool
	profile, _, _, err := repository.GetProfiles(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND user_id = ?",
		FilterValues: []any{userID},
	}, []string{})
	if err != nil || len(profile) == 0 {
		isCreate = true
	}

	var result models.VAProfile
	if isCreate {
		result, statusCode, err = service.CreateProfile(userID, request)
	} else {
		result, statusCode, err = service.UpdateProfile(profile[0].ID.String(), request)
	}
	if err != nil {
		c.AbortWithStatusJSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
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

func GetProfile(c *gin.Context) {
	userID, statusCode, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get current userID",
				Error:   err.Error(),
			},
		)
	}

	profile, _, _, err := repository.GetProfiles(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND user_id = ?",
		FilterValues: []any{userID},
	}, []string{"SocialMedia"})
	if err != nil || len(profile) == 0 {
		if err == nil {
			err = errors.New("profile not found")
		}
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			dto.Response{
				Status:  http.StatusNotFound,
				Message: "Failed to get profile",
				Error:   err.Error(),
			},
		)
	}

	c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  http.StatusOK,
			Message: "Success to update data",
			Data:    profile[0],
		},
	)
}
