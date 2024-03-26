package usecase

import (
	"context"
	"time"

	"promotion/internal/dto"
	"promotion/internal/model"
)

type CampaignUserService struct {
	repo CampaignUserRepo
}

func NewCampaignUserService(repo CampaignUserRepo) *CampaignUserService {
	return &CampaignUserService{repo: repo}
}

func (s *CampaignUserService) CreateCampaignUser(
	ctx context.Context,
	campaignExtID string,
	userExtID string,
	registerDate time.Time,
	firstLoginDate *time.Time,
) (*model.CampaignUser, *model.Campaign, error) {
	var cuser = model.CampaignUser{
		CampaignExtID:    campaignExtID,
		UserExtID:        userExtID,
		RegistrationDate: registerDate,
		FirstLoginDate:   firstLoginDate,
	}
	return s.repo.CreateCampaignUserWithEligibility(ctx, cuser)
}

func (s *CampaignUserService) GetCampaignUsers(ctx context.Context, filter dto.CampaignUserFilter) ([]*model.CampaignUser, error) {
	var m = make(map[string]interface{})
	if len(filter.CampaignExtID) != 0 {
		m["campaign_ext_id"] = filter.CampaignExtID
	}
	if len(filter.UserExtID) != 0 {
		m["user_ext_id"] = filter.UserExtID
	}
	if filter.Vouchered != nil {
		m["vouchered"] = filter.Vouchered
	}
	return s.repo.FindCampaignUsers(ctx, m)
}
