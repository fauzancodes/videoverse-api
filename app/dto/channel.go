package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ChannelRequest struct {
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

func (request ChannelRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Picture, is.URL),
		validation.Field(&request.Name, validation.Required),
	)
}

type SubscribtionRequest struct {
	ChannelID            string `json:"channel_id"`
	NotificationRedirect string `json:"notification_redirect"`
}

func (request SubscribtionRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ChannelID, validation.Required, is.UUID),
	)
}

type UnsubscribeRequest struct {
	NotificationRedirect string `json:"notification_redirect"`
}
