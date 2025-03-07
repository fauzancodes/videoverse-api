package models

import "github.com/google/uuid"

type VAVideoLike struct {
	CustomGormModel
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	VideoID uuid.UUID `json:"video_id" gorm:"type:uuid;column:video_id"`
	User    VAUser    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Video   VAVideo   `json:"video,omitempty" gorm:"foreignKey:VideoID"`
}

func (VAVideoLike) TableName() string {
	return "va_video_likes"
}

type VAVideoDislike struct {
	CustomGormModel
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	VideoID uuid.UUID `json:"video_id" gorm:"type:uuid;column:video_id"`
	User    VAUser    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Video   VAVideo   `json:"video,omitempty" gorm:"foreignKey:VideoID"`
}

func (VAVideoDislike) TableName() string {
	return "va_video_dislikes"
}
