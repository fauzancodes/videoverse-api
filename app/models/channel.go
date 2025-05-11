package models

import "github.com/google/uuid"

type VAChannel struct {
	CustomGormModel
	Name        string           `json:"name" gorm:"type:varchar(255);column:name"`
	Picture     string           `json:"picture" gorm:"type:text;column:picture"`
	Description string           `json:"description" gorm:"type:text;column:description"`
	Location    string           `json:"location" gorm:"type:text;column:location"`
	UserID      uuid.UUID        `json:"user_id" gorm:"type:uuid;column:user_id"`
	Subscribers []VASubscribtion `json:"subscribers,omitempty" gorm:"foreignKey:ChannelID"`
	User        *VAUser          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (VAChannel) TableName() string {
	return "va_channels"
}

type VASubscribtion struct {
	CustomGormModel
	ChannelID    uuid.UUID  `json:"channel_id" gorm:"type:uuid;column:channel_id"`
	SubscriberID uuid.UUID  `json:"subscriber_id" gorm:"type:uuid;column:subscriber_id"`
	Channel      *VAChannel `json:"channel,omitempty" gorm:"foreignKey:ChannelID"`
	Subscriber   *VAUser    `json:"subscriber,omitempty" gorm:"foreignKey:SubscriberID"`
}

func (VASubscribtion) TableName() string {
	return "va_subscribtions"
}
