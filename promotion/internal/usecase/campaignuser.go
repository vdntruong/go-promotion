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

func (s *CampaignUserService) CreateCampaignUser(ctx context.Context, campaignExtID, userExtID string, registerDate time.Time) (*model.CampaignUser, error) {
	var cuser = model.CampaignUser{
		CampaignExtID:    campaignExtID,
		UserExtID:        userExtID,
		RegistrationDate: registerDate,
	}
	return s.repo.CreateCampaignUser(ctx, cuser)
}

func (s *CampaignUserService) DispatchUserLogin(context.Context) error {
	return nil
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
