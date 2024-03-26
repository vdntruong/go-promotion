package voucher

import "github.com/shopspring/decimal"

type CreateVoucher struct {
	CampaignExtID string          `json:"campaignExtId"`
	UserExtID     string          `json:"userExtId"`
	Name          string          `json:"name"`
	Value         decimal.Decimal `json:"value"`
	FixedAmount   bool            `json:"fixedAmount"`
}

type VoucherResponse struct {
	Data struct {
		ExtID         string          `json:"extId"`
		UserExtID     string          `json:"userExtId"`
		CampaignExtID string          `json:"campaignExtId"`
		Name          string          `json:"name"`
		FixedAmount   bool            `json:"fixedAmount"`
		Value         decimal.Decimal `json:"value"`
		IsActive      bool            `json:"isActive"`
		Redeemed      bool            `json:"redeemed"`
	} `json:"data"`
}
