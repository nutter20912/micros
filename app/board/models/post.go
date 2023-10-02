package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `json:"title" gorm:"type:varchar(50);not null"`
	Content string `json:"content" gorm:"type:varchar(500);not null"`
	UserId  string `json:"user_id" gorm:"not null"`

	Comments []Comment `gorm:"foreignKey:PostId"`
}
