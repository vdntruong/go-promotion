package repository

import (
	"context"
	"errors"
	"fmt"

	"promotion/internal/model"
	"promotion/internal/pkg/util"

	"gorm.io/gorm"
)

type CampaignUserRepository struct {
	db *gorm.DB
}

func NewCampaignUserRepository(db *gorm.DB) *CampaignUserRepository {
	return &CampaignUserRepository{db: db}
}

func (c *CampaignUserRepository) CreateCampaignUser(ctx context.Context, cu model.CampaignUser) (*model.CampaignUser, error) {
	if err := c.db.WithContext(ctx).Create(&cu).Error; err != nil {
		if util.IsDuplicateError(err) {
			return nil, model.NewRespError(err, model.WithDescription("duplicate key"))
		}
		return nil, err
	}
	return &cu, nil
}

func (c *CampaignUserRepository) CreateCampaignUserWithEligibility(ctx context.Context, cu model.CampaignUser) (*model.CampaignUser, *model.Campaign, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, nil, err
	}

	var campaign model.Campaign
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("ext_id = ?", cu.CampaignExtID).First(&campaign).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("campaign not found: %w", model.ErrDBNotFound)
		}
		return nil, nil, fmt.Errorf("failed to fetch campaign: %w", err)
	}

	var count int64
	if err := tx.Model(&model.CampaignUser{}).Set("gorm:query_option", "FOR UPDATE").
		Where("campaign_ext_id = ?", cu.CampaignExtID).Count(&count).Error; err != nil {
		tx.Rollback()

		return nil, nil, fmt.Errorf("failed to count campaign users: %w", err)
	}

	if count >= int64(campaign.Eligibility) {
		tx.Rollback()
		return nil, nil, errors.New("maximum number of users registered to the campaign has been reached")
	}

	if err := tx.Create(&cu).Error; err != nil {
		tx.Rollback()
		if util.IsDuplicateError(err) {
			return nil, nil, model.NewRespError(model.ErrDuplicate, model.WithDescription("User has been tracked for this campaign"))
		}
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &cu, &campaign, nil
}

func (c *CampaignUserRepository) FindCampaignUsers(ctx context.Context, filter map[string]interface{}) ([]*model.CampaignUser, error) {
	var result []*model.CampaignUser
	if err := c.db.WithContext(ctx).Where(filter).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
