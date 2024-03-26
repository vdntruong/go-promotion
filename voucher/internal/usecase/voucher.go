package usecase

import (
	"context"
	"errors"
	"net/http"

	"voucher/internal/dto"
	"voucher/internal/model"
)

type VoucherService struct {
	repo      VoucherRepo
	generator VoucherGenerator
}

func NewVoucherService(repo VoucherRepo, generator VoucherGenerator) *VoucherService {
	return &VoucherService{repo: repo, generator: generator}
}

func (s *VoucherService) CreateVoucher(ctx context.Context, v model.Voucher) (*model.Voucher, error) {
	if len(v.Code) == 0 {
		code := s.generator.GetUniqueCode()
		v.Code = code
	}
	return s.repo.CreateVoucher(ctx, v)
}

func (s *VoucherService) RedeemVoucher(ctx context.Context, d dto.RedeemVoucher) (*model.Voucher, error) {
	voucher, found, err := s.repo.FindByExtID(ctx, d.VoucherExtID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, model.NewRespError(
			errors.New("voucher not found"),
			model.WithDescription("Voucher invalid"),
			model.WithStatusCode(http.StatusNoContent),
		)
	}

	if voucher.UserExtID != d.UserExtID {
		return nil, model.NewRespError(
			errors.New("voucher not found"),
			model.WithDescription("User invalid"),
			model.WithStatusCode(http.StatusConflict),
		)
	}

	if voucher.Redeemed {
		return nil, model.NewRespError(
			errors.New("voucher has been redeemed"),
			model.WithDescription("Voucher redeemed"),
			model.WithStatusCode(http.StatusConflict),
		)
	}

	if err := s.repo.RedeemVoucher(ctx, voucher.ExtID); err != nil {
		return nil, model.NewRespError(
			err,
			model.WithDescription("Failed to redeem voucher"),
			model.WithStatusCode(http.StatusConflict),
		)
	}

	result, _, err := s.GetVoucherByExtID(ctx, d.VoucherExtID)
	return result, err
}

func (s *VoucherService) GetVoucherByExtID(ctx context.Context, extID string) (*model.Voucher, bool, error) {
	return s.repo.FindByExtID(ctx, extID)
}

func (s *VoucherService) GetVouchers(ctx context.Context, v dto.VoucherFilter) ([]*model.Voucher, error) {
	var m = make(map[string]interface{})
	if len(v.CampaignExtID) != 0 {
		m["campaign_ext_id"] = v.CampaignExtID
	}
	if len(v.ExtID) != 0 {
		m["ext_id"] = v.ExtID
	}
	if len(v.UserExtID) != 0 {
		m["user_ext_id"] = v.UserExtID
	}
	if len(v.Name) != 0 {
		m["name"] = v.Name
	}
	if v.IsActive != nil {
		m["is_active"] = v.IsActive
	}
	if v.Redeemed != nil {
		m["redeemed"] = v.Redeemed
	}
	return s.repo.FindVouchers(ctx, m)
}
