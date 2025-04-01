package models

type VAUser struct {
	CustomGormModel
	Email           string            `json:"email" gorm:"type:varchar(255);column:email"`
	Password        string            `json:"-" gorm:"type:varchar(255);column:password"`
	IsVerified      bool              `json:"is_verified" gorm:"type:bool;column:is_verified"`
	Profile         *VAProfile        `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Videos          []VAVideo         `json:"videos,omitempty" gorm:"foreignKey:UserID"`
	VideoCategories []VAVideoCategory `json:"video_categories,omitempty" gorm:"foreignKey:UserID"`
	Playlists       []VAPlaylist      `json:"playlists,omitempty" gorm:"foreignKey:UserID"`
	LikedVideos     []VAVideoLike     `json:"liked_videos,omitempty" gorm:"foreignKey:UserID"`
	DislikedVideos  []VAVideoDislike  `json:"disliked_videos,omitempty" gorm:"foreignKey:UserID"`
	Comments        []VAComment       `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

func (VAUser) TableName() string {
	return "va_users"
}
