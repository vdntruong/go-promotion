package model

import (
	"time"

	"gorm.io/gorm"
)

type CampaignUser struct {
	gorm.Model
	CampaignExtID    string    `gorm:"not null;index;uniqueIndex:idx_camp_id_user_id"`
	UserExtID        string    `gorm:"not null;index;uniqueIndex:idx_camp_id_user_id"`
	RegistrationDate time.Time `gorm:"not null"`
	FirstLoginDate   *time.Time
	Vouchered        bool `gorm:"not null;default:false"`
}

func (cu *CampaignUser) TableName() string {
	return "campaign_users"
}
