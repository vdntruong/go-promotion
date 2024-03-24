package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ExtID    string `gorm:"unique;uniqueIndex;not null"`
	Email    string `gorm:"unique;uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Username *string
	FirstLoginDate *time.Time
}

func (c *User) TableName() string {
	return "users"
}

func (c *User) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	c.ExtID = uuid
	return nil
}
