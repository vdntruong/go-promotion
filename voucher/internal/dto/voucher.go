package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateVoucher struct {
	CampaignExtID string `json:"campaignExtId" binding:"required"`
	UserExtID     string `json:"userExtId" binding:"required"`
	Name          string `json:"name" binding:"required"`

	FixedAmount bool            `json:"fixedAmount"`
	Value       decimal.Decimal `json:"value" binding:"required"`
}

type RedeemVoucher struct {
	RequestID    string `json:"requestID" binding:"required"`
	UserExtID    string `json:"userExtId" binding:"required"`
	VoucherExtID string `json:"voucherExtID" binding:"required"`

	Cost  decimal.Decimal `json:"cost" binding:"required"`
	Price decimal.Decimal `json:"price" binding:"required"`
}

type VoucherFilter struct {
	CampaignExtID string `form:"campaignExtId"`
	UserExtID     string `form:"userExtId"`
	Name          string `form:"name"`
	ExtID         string `form:"extId"`

	IsActive *bool `form:"isActive"`
	Redeemed *bool `form:"redeemed"`
}

type Voucher struct {
	ExtID         string `json:"extId"`
	UserExtID     string `json:"userExtId"`
	CampaignExtID string `json:"campaignExtId"`
	Name          string `json:"name"`

	FixedAmount bool            `json:"fixedAmount"`
	Value       decimal.Decimal `json:"value"`

	IsActive bool       `json:"isActive"`
	Redeemed bool       `json:"redeemed"`
	RedeemAt *time.Time `json:"redeemAt,omitempty"`
}
