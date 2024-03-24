package dto

import "time"

type Campaign struct {
	ExtID       string    `json:"extID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	IsActive    bool      `json:"isActive"`
	Eligibility int       `json:"eligibility"`
}

type CreateCampaign struct {
	Name        string    `form:"name" json:"name" binding:"required"`
	Description string    `form:"description" json:"description" binding:"required"`
	StartDate   time.Time `form:"startDate" json:"startDate" binding:"required"`
	EndDate     time.Time `form:"endDate" json:"endDate" binding:"required"`
}
