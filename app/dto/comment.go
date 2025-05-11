package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CommentRequest struct {
	ParentID             string `json:"parent_id"`
	VideoID              string `json:"video_id"`
	Content              string `json:"content"`
	NotificationRedirect string `json:"notification_redirect"`
}

func (request CommentRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ParentID, is.UUID),
		validation.Field(&request.VideoID, validation.Required, is.UUID),
		validation.Field(&request.Content, validation.Required),
	)
}

type CommentUpdateRequest struct {
	Content string `json:"content"`
}

func (request CommentUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Content, validation.Required),
	)
}
