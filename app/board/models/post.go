package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(50);not null"`
	Content string `gorm:"type:varchar(500);not null"`
	UserId  uint   `gorm:"not null"`
}
