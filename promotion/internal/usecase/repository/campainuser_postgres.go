package repository

import (
	"context"
	"time"

	"promotion/internal/model"
	"promotion/internal/pkg/util"

	"gorm.io/gorm"
)

type CampaignUserRepository struct {
	db *gorm.DB
}

func NewCampaignUserRepository(db *gorm.DB) *CampaignUserRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
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

func (c *CampaignUserRepository) FindCampaignUsers(ctx context.Context, filter map[string]interface{}) ([]*model.CampaignUser, error) {
	var result []*model.CampaignUser
	if err := c.db.WithContext(ctx).Where(filter).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
