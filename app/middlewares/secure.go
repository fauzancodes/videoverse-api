package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/config"
	"github.com/fauzancodes/videoverse-api/app/dto"
	"github.com/fauzancodes/videoverse-api/app/models"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	strip "github.com/grokify/html-strip-tags-go"
)

// Secure Middleware
func Secure() gin.HandlerFunc {
	return secure.New(secure.Config{
		// SSLRedirect:   false,
		IsDevelopment: false,
		// STSSeconds:            315360000,
		// STSIncludeSubdomains:  true,
		// STSPreload:            true,
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: "default-src 'self'",
		IENoOpen: true,
		// SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func StripHTMLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := c.Request.URL.Query()
		for key, values := range queryParams {
			for i, value := range values {
				sanitizedValue := template.HTMLEscapeString(value)
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "=", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "<", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, ">", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "*", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " AND ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " OR ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " and ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " or ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " || ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " && ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "'", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "&#39;", "")
				values[i] = strip.StripTags(sanitizedValue)
			}
			queryParams[key] = values
		}
		c.Request.URL.RawQuery = queryParams.Encode()
		c.Next()
	}
}

func CheckAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.LoadConfig()

		if !conf.EnableAPIKey {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			fmt.Println("Failed to check API key: no API key in header")
			c.JSON(http.StatusForbidden, dto.Response{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
			})
			c.Abort()
			return
		}

		if apiKey == conf.SpecialApiKey {
			c.Next()
			return
		}

		// Decode API Key
		secretKey, receivedHMAC, err := DecodeAPIKeyBase64(apiKey)
		if err != nil {
			fmt.Println("Failed to decode API key: ", err.Error())
			c.JSON(http.StatusForbidden, dto.Response{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
			})
			c.Abort()
			return
		}

		// Verify HMAC
		secret := conf.HMACKey
		hmacVerified, expectedHMAC, err := VerifyAPIKeyHMAC(secretKey, receivedHMAC, secret)
		if err != nil || !hmacVerified {
			fmt.Println("Failed to verify API key: ", err.Error())
			c.JSON(http.StatusForbidden, dto.Response{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
			})
			c.Abort()
			return
		}

		// Check if API Key already used
		var usedApiKey models.VAUsedApiKey
		err = config.DB.Debug().Where("secret_key = ?", secretKey).First(&usedApiKey).Error
		if err == nil && usedApiKey.ID != uuid.Nil {
			fmt.Println("Failed to check API key: API key already used")
			c.JSON(http.StatusForbidden, dto.Response{
				Status:  http.StatusForbidden,
				Message: "Forbidden",
			})
			c.Abort()
			return
		}

		// Save Used API Key
		err = config.DB.Create(&models.VAUsedApiKey{
			SecretKey:    secretKey,
			Base64Key:    apiKey,
			ReceivedHMAC: receivedHMAC,
			ExpectedHMAC: expectedHMAC,
		}).Error
		if err != nil {
			fmt.Println("Failed to save API key: ", err.Error())
		}

		c.Next()
	}
}

func ComputeAPIKeyHMAC(secretKey, secret string) (response string, err error) {
	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write([]byte(secretKey))
	response = hex.EncodeToString(h.Sum(nil))

	return
}

func DecodeAPIKeyBase64(encodedKey string) (secretKey string, hmacSignature string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		err = errors.New("invalid encoded key")
		return
	}

	secretKey = parts[0]
	hmacSignature = parts[1]

	return
}

func VerifyAPIKeyHMAC(secretKey, receivedHMAC, secret string) (response bool, expectedHMAC string, err error) {
	expectedHMAC, err = ComputeAPIKeyHMAC(secretKey, secret)
	if err != nil {
		return
	}
	response = hmac.Equal([]byte(expectedHMAC), []byte(receivedHMAC))

	return
}
