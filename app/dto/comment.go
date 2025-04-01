package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CommentRequest struct {
	ParentID string `json:"parent_id" gorm:"type:uuid;column:parent_id"`
	VideoID  string `json:"video_id" gorm:"type:uuid;column:video_id"`
	Content  string `json:"content" gorm:"type:text;column:content"`
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
	Content string `json:"content" gorm:"type:text;column:content"`
}

func (request CommentUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Content, validation.Required),
	)
}
