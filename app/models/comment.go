package models

import "github.com/google/uuid"

type VAComment struct {
	CustomGormModel
	ParentID *uuid.UUID  `json:"parent_id" gorm:"type:uuid;column:parent_id"`
	UserID   uuid.UUID   `json:"user_id" gorm:"type:uuid;column:user_id"`
	VideoID  uuid.UUID   `json:"video_id" gorm:"type:uuid;column:video_id"`
	Content  string      `json:"content" gorm:"type:text;column:content"`
	Parent   *VAComment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies  []VAComment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	User     *VAUser     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Video    *VAVideo    `json:"video,omitempty" gorm:"foreignKey:VideoID"`
}

func (VAComment) TableName() string {
	return "va_comments"
}
