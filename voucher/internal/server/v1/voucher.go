package v1

import (
	"net/http"

	"voucher/internal/dto"
	"voucher/internal/model"
	"voucher/internal/server/v1/resp"
	"voucher/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type voucherRoutes struct {
	c usecase.VoucherUsecase
}

func newVoucherRoutes(handler *gin.RouterGroup, c usecase.VoucherUsecase) {
	r := &voucherRoutes{c}
	h := handler.Group("/vouchers")
	{
		h.POST("/", r.CreateVoucher)
		h.GET("/", r.GetVouchers)
		h.POST("/redeem", r.RedeemVoucher)
	}
}

func (r *voucherRoutes) CreateVoucher(ctx *gin.Context) {
	var req dto.CreateVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	var obj model.Voucher
	copier.Copy(&obj, &req)

	result, err := r.c.CreateVoucher(ctx, obj)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var res dto.Voucher
	copier.Copy(&res, &result)
	resp.RespondDataWithCode(ctx, res, http.StatusCreated)
}

func (r *voucherRoutes) RedeemVoucher(ctx *gin.Context) {
	var req dto.RedeemVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	result, err := r.c.RedeemVoucher(ctx, req)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var res dto.Voucher
	copier.Copy(&res, &result)
	resp.RespondDataWithCode(ctx, res, http.StatusOK)
}


func (r *voucherRoutes) GetVouchers(ctx *gin.Context) {
	var filter dto.VoucherFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		resp.RespondError(ctx, err)
		return
	}

	result, err := r.c.GetVouchers(ctx.Request.Context(), filter)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var res []*dto.Voucher
	copier.Copy(&res, &result)
	resp.RespondData(ctx, res)
}
