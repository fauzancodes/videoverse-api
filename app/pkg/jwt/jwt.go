package jwt

import (
	"errors"
	"fmt"

	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = config.LoadConfig().SecretKey

func GenerateToken(claims *jwt.MapClaims) (string, error) {
	//encode header and payload into token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign the token with a secret key
	//then combine the encoded header, payload and signature
	webtoken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		err = fmt.Errorf("failed to generate token: %s", err.Error())
		return "", err
	}

	return webtoken, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	//verify payload and signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		err = fmt.Errorf("failed to verify token: %s", err.Error())
		return token, err
	}

	return token, nil
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	//vefify token
	token, err := VerifyToken(tokenString)
	claims, isOk := token.Claims.(jwt.MapClaims)
	if !isOk {
		err = errors.New("failed to infer token claims")
		return claims, err
	}
	if err != nil {
		err = fmt.Errorf("failed to decode token: %s", err.Error())
		return claims, err
	}

	if token.Valid {
		return claims, nil
	}

	return claims, fmt.Errorf("invalid token")
}
