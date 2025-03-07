package dto

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type PlaylistRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      bool     `json:"status"`
	Visibility  string   `json:"visibility"`
	VideoIDs    []string `json:"video_ids"`
}

func (request PlaylistRequest) Validate() error {
	visibilities := []string{"public", "unlisted", "private"}
	var visibilityAccepted bool
	for _, visibility := range visibilities {
		if strings.EqualFold(request.Visibility, visibility) {
			visibilityAccepted = true
		}
	}
	if !visibilityAccepted {
		return fmt.Errorf("allowed visibilities: %s", strings.Join(visibilities, ", "))
	}

	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Title, validation.Required),
	)
}
