package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

func BuildPreload(db *gorm.DB, fields []string) *gorm.DB {
	if len(fields) > 0 {
		for _, field := range fields {
			db = db.Preload(field)
		}
	}

	return db
}

func GetBuildPreloadFields(c echo.Context) (fields []string) {
	raw := c.QueryParam("preload_fields")

	if raw != "" {
		fields = strings.Split(raw, ",")
	}

	return
}

func GenerateRandomNumber(length int) string {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}
	rand.Seed(uint64(time.Now().In(location).UnixNano()))
	charset := "0123456789"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(randomBytes)
	return randomString
}

func GetBaseUrl(c *gin.Context) (response string) {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	response = fmt.Sprintf("%v://%v", scheme, c.Request.Host)

	return
}

func RoundFloat(number float64) (result int) {
	return int(math.Round(number))
}

func GetCurrentUserID(c *gin.Context) (userID string, statusCode int, err error) {
	authContext, exist := c.Get("currentUser")
	if !exist {
		err = errors.New("no current user data in context")
		statusCode = http.StatusInternalServerError

		return
	}

	userID = authContext.(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	return
}

func IsValidUUID(id string) bool {
	if id == "" {
		return true
	}

	_, err := uuid.Parse(id)

	return err == nil
}

func QueryParamUUID(c *gin.Context, key string) (id string, statusCode int, err error) {
	id = c.Query(key)
	validUUID := IsValidUUID(id)
	if !validUUID {
		err = fmt.Errorf("%s: %s, is not valid UUID", key, id)
		statusCode = http.StatusBadRequest
		return
	}

	statusCode = http.StatusOK
	return
}

func ParamUUID(c *gin.Context, key string) (id string, statusCode int, err error) {
	id = c.Param(key)
	validUUID := IsValidUUID(id)
	if !validUUID {
		err = fmt.Errorf("%s: %s, is not valid UUID", key, id)
		statusCode = http.StatusBadRequest
		return
	}

	statusCode = http.StatusOK
	return
}

func GetTimeLocationAsiaJakarta() (location *time.Location) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}

	return
}
