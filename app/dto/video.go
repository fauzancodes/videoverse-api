package dto

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type VideoRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	CategoryID   string   `json:"category_id"`
	VideoUrl     string   `json:"video_url"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	Visibility   string   `json:"visibility"`
	Tags         []string `json:"tags"`
	Status       bool     `json:"status"`
}

func (request VideoRequest) Validate() error {
	visibilities := []string{"public", "unlisted", "private"}
	var visibilityAccepted bool
	for _, visibility := range visibilities {
		if strings.EqualFold(request.Visibility, visibility) {
			visibilityAccepted = true
		}
	}
	if !visibilityAccepted {
		return fmt.Errorf("allowed visibilities: %s", strings.Join(visibilities, ", "))
	}

	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Title, validation.Required),
		validation.Field(&request.VideoUrl, validation.Required, is.URL),
		validation.Field(&request.Visibility, validation.Required),
		validation.Field(&request.ThumbnailUrl, is.URL),
		validation.Field(&request.CategoryID, is.UUID),
	)
}

type VideoCategoryRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func (request VideoCategoryRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Title, validation.Required),
	)
}
