package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type VideoLikeRequest struct {
	VideoID string `json:"video_id"`
}

func (request VideoLikeRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.VideoID, validation.Required, is.UUID),
	)
}

type VideoDislikeRequest struct {
	VideoID string `json:"video_id"`
}

func (request VideoDislikeRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.VideoID, validation.Required, is.UUID),
	)
}
