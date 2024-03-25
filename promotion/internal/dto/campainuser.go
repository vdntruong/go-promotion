package dto

import "time"

type CampaignUser struct {
	CampaignExtID    string     `json:"campaignExtId"`
	UserExtID        string     `json:"userExtId"`
	RegistrationDate time.Time  `json:"registrationDate"`
	FirstLoginDate   *time.Time `json:"firstLoginDate,omitempty"`
	Vouchered        bool       `json:"vouchered"`
}

type CreateCampaignUser struct {
	CampaignExtID string    `form:"campaignExtId" json:"campaignExtId" binding:"required"`
	UserExtID     string    `form:"userExtId" json:"userExtId" binding:"required"`
	RegisterDate  time.Time `form:"registerDate" json:"registerDate" binding:"required"`
}

type CampaignUserFilter struct {
	CampaignExtID string `form:"campaignExtId"`
	UserExtID     string `form:"userExtId"`
	Vouchered     *bool  `form:"vouchered,omitempty"`
}
