package models

type VAUser struct {
	CustomGormModel
	Email      string     `json:"email" gorm:"type:varchar(255);column:email"`
	Password   string     `json:"-" gorm:"type:varchar(255);column:password"`
	IsVerified bool       `json:"is_verified" gorm:"type:bool;column:is_verified"`
	Profile    *VAProfile `json:"profile,omitempty" gorm:"foreignKey:UserID"`
}

func (VAUser) TableName() string {
	return "va_users"
}
