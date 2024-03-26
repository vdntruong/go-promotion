package usecase

import (
	"context"

	"voucher/internal/dto"
	"voucher/internal/model"
)

type VoucherUsecase interface {
	RedeemVoucher(context.Context, dto.RedeemVoucher) (*model.Voucher, error)
	CreateVoucher(context.Context, model.Voucher) (*model.Voucher, error)
	GetVoucherByExtID(context.Context, string) (*model.Voucher, bool, error)
	GetVouchers(context.Context, dto.VoucherFilter) ([]*model.Voucher, error)
}

type VoucherRepo interface {
	CreateVoucher(context.Context, model.Voucher) (*model.Voucher, error)
	FindVouchers(context.Context, map[string]interface{}) ([]*model.Voucher, error)
	FindByExtID(context.Context, string) (*model.Voucher, bool, error)
	FindByVoucherCode(context.Context, string) (*model.Voucher, bool, error)
	FindByCampaignAndUserExtID(context.Context, string, string) (*model.Voucher, bool, error)
	RedeemVoucher(context.Context, string) error
}

type VoucherGenerator interface {
	GetUniqueCode() string
}
