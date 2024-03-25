package repository

import (
	"context"

	"promotion/internal/model"
	"promotion/internal/pkg/util"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (c *CampaignRepository) CreateCampaign(ctx context.Context, cp model.Campaign) (*model.Campaign, error) {
	if err := c.db.WithContext(ctx).Create(&cp).Error; err != nil {
		if util.IsDuplicateError(err) {
			return nil, model.NewRespError(err, model.WithDescription("duplicate key"))
		}
		return nil, err
	}
	return &cp, nil
}

func (c *CampaignRepository) FindByID(context.Context, int) (*model.Campaign, bool, error) {
	return nil, false, nil
}

func (c *CampaignRepository) FindByExtID(context.Context, string) (*model.Campaign, bool, error) {
	return nil, false, nil
}

func (c *CampaignRepository) FindCampaigns(ctx context.Context) ([]*model.Campaign, error) {
	var result []*model.Campaign
	if err := c.db.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
