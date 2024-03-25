package dto

type User struct {
	ExtID    string  `json:"extID"`
	Email    string  `json:"email"`
	Username *string `json:"username,omitempty"`
}

type SignUpUser struct {
	Email         string  `json:"email" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	CampaignExtID *string `json:"campaignExtID,omitempty"`
}

type SignInUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
