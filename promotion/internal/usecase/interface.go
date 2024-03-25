package usecase

import (
	"context"
	"time"

	"promotion/internal/dto"
	"promotion/internal/model"
)

type CampaignUsecase interface {
	CreateCampaign(context.Context, model.Campaign) (*model.Campaign, error)
	GetCampaigns(context.Context) ([]*model.Campaign, error)
	GetCampaignByID(context.Context, int) (*model.Campaign, bool, error)
	GetCampaignByExtID(context.Context, string) (*model.Campaign, bool, error)
	// DeactivateCampaign(ctx context.Context, id int) (*model.Campaign, error)
}

type CampaignRepo interface {
	CreateCampaign(context.Context, model.Campaign) (*model.Campaign, error)
	FindCampaigns(context.Context) ([]*model.Campaign, error)
	FindByID(context.Context, int) (*model.Campaign, bool, error)
	FindByExtID(context.Context, string) (*model.Campaign, bool, error)
}

type CampaignUserUsecase interface {
	CreateCampaignUser(ctx context.Context, campaignExtID, userExtID string, registerDate time.Time) (*model.CampaignUser, *model.Campaign, error)
	GetCampaignUsers(ctx context.Context, filter dto.CampaignUserFilter) ([]*model.CampaignUser, error)
}

type CampaignUserRepo interface {
	CreateCampaignUser(context.Context, model.CampaignUser) (*model.CampaignUser, error)
	CreateCampaignUserWithEligibility(ctx context.Context, cu model.CampaignUser) (*model.CampaignUser, *model.Campaign, error)
	FindCampaignUsers(context.Context, map[string]interface{}) ([]*model.CampaignUser, error)
}
