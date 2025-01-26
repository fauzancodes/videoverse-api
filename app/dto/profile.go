package dto

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ProfileRequest struct {
	Firstname   string               `json:"firstname"`
	Lastname    string               `json:"lastname"`
	Gender      string               `json:"gender"`
	Picture     string               `json:"picture"`
	Description string               `json:"description"`
	Location    string               `json:"location"`
	SocialMedia []SocialMediaRequest `json:"social_media"`
}

func (request ProfileRequest) Validate() error {
	genders := []string{"male", "female"}
	var genderAccepted bool
	for _, gender := range genders {
		if strings.EqualFold(request.Gender, gender) {
			genderAccepted = true
		}
	}
	if !genderAccepted {
		return errors.New("allowed genders: " + strings.Join(genders, ", "))
	}

	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Picture, is.URL),
		validation.Field(&request.Firstname, validation.Required),
		validation.Field(&request.Lastname, validation.Required),
	)
}

type SocialMediaRequest struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func (request SocialMediaRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Link, is.URL),
	)
}
