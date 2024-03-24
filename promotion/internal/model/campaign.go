package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campaign struct {
    gorm.Model
    ExtID        string    `gorm:"unique;uniqueIndex;not null"`
    Name         string    `gorm:"unique;not null"`
    Description  string    `gorm:"not null"`
    StartDate    time.Time `gorm:"not null"`
    EndDate      time.Time `gorm:"not null"`
    IsActive     bool      `gorm:"not null;default:true"`
    Eligibility  int       `gorm:"not null;default:100"`
}

func (c *Campaign) TableName() string {
    return "campaigns"
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (error) {
    uuid := uuid.New().String()
    c.ExtID = uuid
    return nil
}