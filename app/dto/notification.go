package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type NotificationRequest struct {
	Redirect string `json:"redirect"`
	Content  string `json:"content"`
}

func (request NotificationRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Content, validation.Required),
		validation.Field(&request.Redirect, validation.Required),
	)
}
