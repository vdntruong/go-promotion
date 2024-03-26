package repository

import (
	"context"
	"errors"
	"time"

	"voucher/internal/model"
	"voucher/internal/pkg/util"

	"gorm.io/gorm"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

func (c *VoucherRepository) CreateVoucher(ctx context.Context, v model.Voucher) (*model.Voucher, error) {
	if err := c.db.WithContext(ctx).Create(&v).Error; err != nil {
		if util.IsDuplicateError(err) {
			voucher, found, findErr := c.FindByCampaignAndUserExtID(ctx, v.CampaignExtID, v.UserExtID)
			if findErr != nil {
				return nil, err
			}
			if !found {
				return nil, model.NewRespError(err, model.WithDescription("duplicate key"))
			}
			return voucher, nil
		}
		return nil, err
	}
	return &v, nil
}

func (c *VoucherRepository) FindByExtID(ctx context.Context, extID string) (*model.Voucher, bool, error) {
	var result *model.Voucher
	if err := c.db.WithContext(ctx).First(&result, "ext_id = ?", extID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return result, true, nil
}

func (c *VoucherRepository) FindByVoucherCode(ctx context.Context, code string) (*model.Voucher, bool, error) {
	var result *model.Voucher
	if err := c.db.WithContext(ctx).First(&result, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return result, true, nil
}

func (c *VoucherRepository) FindByCampaignAndUserExtID(
	ctx context.Context,
	campaignExtID string,
	userExtID string,
) (*model.Voucher, bool, error) {
	var result *model.Voucher
	if err := c.db.WithContext(ctx).
		First(&result, "campaign_ext_id = ? AND user_ext_id = ?", campaignExtID, userExtID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return result, true, nil
}

func (c *VoucherRepository) FindVouchers(ctx context.Context, filter map[string]interface{}) ([]*model.Voucher, error) {
	var result []*model.Voucher
	if err := c.db.WithContext(ctx).Where(filter).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *VoucherRepository) IsUnique(code string) bool {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancelFunc()

	_, found, err := s.FindByVoucherCode(ctx, code)
	if err != nil || found {
		return false
	}

	return true
}

func (s *VoucherRepository) RedeemVoucher(ctx context.Context, extID string) (error) {
	if err := s.db.WithContext(ctx).Model(&model.Voucher{}).
		Where("ext_id = ? AND redeemed = ? AND is_active = ?", extID, false, true).
		Update("redeemed", true).
		Update("redeem_at", util.Ptr(time.Now())).Error; err != nil {
		return  err
	}
	return  nil
}
