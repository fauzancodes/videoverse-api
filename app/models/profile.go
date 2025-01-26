package models

import "github.com/google/uuid"

type VAProfile struct {
	CustomGormModel
	Firstname   string          `json:"firstname" gorm:"type:varchar(255);column:firstname"`
	Lastname    string          `json:"lastname" gorm:"type:varchar(255);column:lastname"`
	Gender      string          `json:"gender" gorm:"type:varchar(50);column:gender"`
	Picture     string          `json:"picture" gorm:"type:text;column:picture"`
	Description string          `json:"description" gorm:"type:text;column:description"`
	Location    string          `json:"location" gorm:"type:text;column:location"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;column:user_id"`
	SocialMedia []VASocialMedia `json:"social_media,omitempty" gorm:"foreignKey:ProfileID"`
}

func (VAProfile) TableName() string {
	return "va_profiles"
}

type VASocialMedia struct {
	CustomGormModel
	Name      string    `json:"name" gorm:"type:varchar(50);column:name"`
	Link      string    `json:"link" gorm:"type:text;column:link"`
	ProfileID uuid.UUID `json:"profile_id" gorm:"type:uuid;column:profile_id"`
}

func (VASocialMedia) TableName() string {
	return "va_social_medias"
}
