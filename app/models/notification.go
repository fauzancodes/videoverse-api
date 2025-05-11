package models

import "github.com/google/uuid"

type VANotification struct {
	CustomGormModel
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	Content  string    `json:"content" gorm:"type:text;column:content"`
	IsRead   bool      `json:"is_read" gorm:"type:bool;column:is_read"`
	Redirect string    `json:"redirect" gorm:"type:varchar(255);column:redirect"`
	User     *VAUser   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (VANotification) TableName() string {
	return "va_notifications"
}
