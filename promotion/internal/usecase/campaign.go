package usecase

import (
	"context"

	"promotion/internal/model"
)

type CampaignService struct {
	repo CampaignRepo
}

func NewCampaignService(r CampaignRepo) *CampaignService {
	return &CampaignService{
		repo: r,
	}
}

func (c *CampaignService) CreateCampaign(ctx context.Context, cp model.Campaign) (*model.Campaign, error) {
	result, err := c.repo.CreateCampaign(ctx, cp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CampaignService) GetCampaignByID(ctx context.Context, id int) (*model.Campaign, bool, error) {
	return c.repo.FindByID(ctx, id)
}

func (c *CampaignService) GetCampaignByExtID(ctx context.Context, extID string) (*model.Campaign, bool, error) {
	return c.repo.FindByExtID(ctx, extID)
}

func (c *CampaignService) GetCampaigns(ctx context.Context) ([]*model.Campaign, error) {
	return c.repo.FindCampaigns(ctx)
}
