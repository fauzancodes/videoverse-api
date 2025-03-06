package models

import (
	"github.com/google/uuid"
)

type VAVideo struct {
	CustomGormModel
	Title        string           `json:"title" gorm:"type:varchar(255);column:title"`
	Description  string           `json:"description" gorm:"type:text;column:description"`
	CategoryID   *uuid.UUID       `json:"category_id" gorm:"type:uuid;column:category_id"`
	VideoUrl     string           `json:"video_url" gorm:"type:text;column:video_url"`
	ThumbnailUrl string           `json:"thumbnail_url" gorm:"type:text;column:thumbnail_url"`
	Visibility   string           `json:"visibility" gorm:"type:varchar(50);column:visibility"`
	Tags         string           `json:"tags" gorm:"type:varchar(255);column:tags"`
	Status       bool             `json:"status" gorm:"type:bool;column:status"`
	UserID       *uuid.UUID       `json:"user_id" gorm:"type:uuid;column:user_id"`
	User         *VAUser          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category     *VAVideoCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Playlists    []VAPlaylist     `json:"playlists,omitempty" gorm:"many2many:va_video_playlists"`
}

func (VAVideo) TableName() string {
	return "va_videos"
}

type VAVideoCategory struct {
	CustomGormModel
	Title       string     `json:"title" gorm:"type:varchar(255);column:title"`
	Description string     `json:"description" gorm:"type:text;column:description"`
	Status      bool       `json:"status" gorm:"type:bool;column:status"`
	UserID      *uuid.UUID `json:"-" gorm:"type:uuid;column:user_id"`
	User        *VAUser    `json:"-" gorm:"foreignKey:UserID"`
	Videos      []VAVideo  `json:"videos,omitempty" gorm:"foreignKey:CategoryID"`
}

func (VAVideoCategory) TableName() string {
	return "va_video_categories"
}

type VAPlaylist struct {
	CustomGormModel
	Title       string     `json:"title" gorm:"type:varchar(255);column:title"`
	Description string     `json:"description" gorm:"type:text;column:description"`
	Status      bool       `json:"status" gorm:"type:bool;column:status"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	User        *VAUser    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Videos      []VAVideo  `json:"videos,omitempty" gorm:"many2many:va_video_playlists"`
}

func (VAPlaylist) TableName() string {
	return "va_playlists"
}
