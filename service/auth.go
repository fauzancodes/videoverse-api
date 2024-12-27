package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/fauzancodes/videoverse-api/app/pkg/bcrypt"
	webToken "github.com/fauzancodes/videoverse-api/app/pkg/jwt"
	"github.com/fauzancodes/videoverse-api/app/pkg/smtp"
	"github.com/fauzancodes/videoverse-api/app/repository"
	"github.com/golang-jwt/jwt/v5"
)

func SendEmailVerification(user models.VAUser, successUrl, failedUrl, appUrl string) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["successUrl"] = successUrl
	claims["failedUrl"] = failedUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println(err.Error())
		return
	}

	verificationUrl := fmt.Sprintf("%v/v1/auth/email-verification/%v", appUrl, token)

	fill := dto.EmailVerfication{
		Name:            user.Email,
		AppUrl:          appUrl,
		VerificationUrl: verificationUrl,
	}

	smtp.SendEmail("email-verification", "", user.Email, "Email Verification", "", fill)
}

func VerifyUser(token string) (user models.VAUser, successUrl, failedUrl string, err error) {
	if token == "" {
		err = errors.New("no jwt token provided")
		return
	}

	claims, err := webToken.DecodeToken(token)
	if err != nil {
		log.Println(err.Error())
		return
	}

	successUrl = claims["successUrl"].(string)
	failedUrl = claims["failedUrl"].(string)

	userID := claims["id"].(string)
	user, _, err = GetUserByID(userID, []string{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	user.IsVerified = true
	user, err = repository.UpdateUser(user)
	if err != nil {
		log.Println("Failed to update user: " + err.Error())
		return
	}

	return
}

func SendResetPasswordRequest(user models.VAUser, redirectUrl, appUrl string) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["redirectUrl"] = redirectUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var resetPasswordUrl string
	if redirectUrl != "" {
		resetPasswordUrl = fmt.Sprintf("%v/%v", redirectUrl, token)
	} else {
		resetPasswordUrl = fmt.Sprintf("%v/v1/auth/reset-password/instruction/%v", appUrl, token)
	}

	fill := dto.ResetPassword{
		Name:             user.Email,
		AppUrl:           appUrl,
		ResetPasswordUrl: resetPasswordUrl,
	}

	smtp.SendEmail("reset-password", "", user.Email, "Reset Your Password", "", fill)
}

func ResetPassword(request dto.ResetPasswordRequest) (user models.VAUser, statusCode int, err error) {
	claims, err := webToken.DecodeToken(request.Token)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	userID := claims["id"].(string)
	user, statusCode, err = GetUserByID(userID, []string{})
	if err != nil {
		return
	}

	user.Password = bcrypt.HashPassword(request.NewPassword)
	user, err = repository.UpdateUser(user)
	if err != nil {
		log.Println("Failed to update user: ", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	return
}
