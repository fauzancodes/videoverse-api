package dto

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EmailVerfication struct {
	Name            string
	VerificationUrl string
	AppUrl          string
}

type ResetPassword struct {
	Name             string
	ResetPasswordUrl string
	AppUrl           string
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

func (request ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.NewPassword, validation.Required),
		validation.Field(&request.Token, validation.Required),
	)
}

type RegisterRequest struct {
	Firstname              string `json:"firstname"`
	Lastname               string `json:"lastname"`
	Gender                 string `json:"gender"`
	Email                  string `json:"email"`
	Password               string `json:"password"`
	SuccessVerificationUrl string `json:"success_verification_url"`
	FailedVerificationUrl  string `json:"failed_verification_url"`
}

func (request RegisterRequest) Validate() error {
	genders := []string{"male", "female"}
	var genderAccepted bool
	for _, gender := range genders {
		if strings.EqualFold(request.Gender, gender) {
			genderAccepted = true
		}
	}
	if !genderAccepted {
		return fmt.Errorf("allowed genders: %s", strings.Join(genders, ", "))
	}

	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
		validation.Field(&request.SuccessVerificationUrl, is.URL),
		validation.Field(&request.FailedVerificationUrl, is.URL),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request LoginRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
	)
}

type ResendEmailVerification struct {
	Email                  string `json:"email"`
	SuccessVerificationUrl string `json:"success_verification_url"`
	FailedVerificationUrl  string `json:"failed_verification_url"`
}

func (request ResendEmailVerification) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.SuccessVerificationUrl, is.URL),
		validation.Field(&request.FailedVerificationUrl, is.URL),
	)
}

type SendForgotPasswordRequest struct {
	Email       string `json:"email"`
	RedirectUrl string `json:"redirect_url"`
}

func (request SendForgotPasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.RedirectUrl, is.URL),
	)
}
