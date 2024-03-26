package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Voucher struct {
	gorm.Model
	ExtID         string `gorm:"uniqueIndex;not null"`
	Name          string `gorm:"not null"`
	UserExtID     string `gorm:"uniqueIndex:idx_camp_id_user_id"`
	CampaignExtID string `gorm:"uniqueIndex:idx_camp_id_user_id"`

	Code        string          `gorm:"uniqueIndex;not null"`
	FixedAmount bool            `gorm:"not null"`
	Value       decimal.Decimal `gorm:"type:decimal(10,2);not null"`

	IsActive bool `gorm:"not null;default:true"`
	Redeemed bool `gorm:"not null;default:false"`
	RedeemAt *time.Time
}

func (c *Voucher) TableName() string {
	return "vouchers"
}

func (c *Voucher) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	c.ExtID = uuid
	return nil
}
