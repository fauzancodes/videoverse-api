package models

import "github.com/google/uuid"

type VAPlaylist struct {
	CustomGormModel
	Title       string     `json:"title" gorm:"type:varchar(255);column:title"`
	Description string     `json:"description" gorm:"type:text;column:description"`
	Status      bool       `json:"status" gorm:"type:bool;column:status"`
	Visibility  string     `json:"visibility" gorm:"type:varchar(50);column:visibility"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	User        *VAUser    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Videos      []VAVideo  `json:"videos,omitempty" gorm:"many2many:va_video_playlists"`
}

func (VAPlaylist) TableName() string {
	return "va_playlists"
}
