package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId  string `json:"user_id" gorm:"not null"`
	PostId  string `json:"post_id" gorm:"not null"`
	Content string `json:"content" gorm:"type:varchar(500);not null"`

	Post Post
}
