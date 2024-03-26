package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"promotion/internal/model"
	"promotion/internal/usecase"

	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
)

type UserTaskHandler struct {
	uc        usecase.CampaignUserUsecase
	campUc    usecase.CampaignUsecase
	voucherUc usecase.VoucherUsecase
}

func NewUserTaskHandler(c usecase.CampaignUserUsecase, campUc usecase.CampaignUsecase, voucherUc usecase.VoucherUsecase) *UserTaskHandler {
	return &UserTaskHandler{uc: c, campUc: campUc, voucherUc: voucherUc}
}

func (u *UserTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p UserFirstLoginPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return u.handleUserFirstTimeLogin(ctx, p)
}

func (u *UserTaskHandler) handleUserFirstTimeLogin(ctx context.Context, p UserFirstLoginPayload) error {
	campU, camp, err := u.uc.CreateCampaignUser(ctx, p.CampaignExtID, p.UserExtID, p.RegisterTime, &p.LoginDateTime)
	if err != nil {
		if errors.Is(err, model.ErrDBNotFound) || errors.Is(err, model.ErrDuplicate) {
			return fmt.Errorf("failed to add campaign user record: %v: %w", err, asynq.SkipRetry)
		}
		return fmt.Errorf("failed to add campaign user record: %v", err)
	}

	if campU.FirstLoginDate.After(camp.StartDate) && campU.FirstLoginDate.Before(camp.EndDate) {
		retry, err := u.voucherUc.CreateVoucher(
			ctx,
			camp.Name,
			campU.CampaignExtID,
			campU.UserExtID,
			decimal.NewFromFloat(0.3),
		)
		if err != nil {
			if retry {
				return fmt.Errorf("failed to generate voucher record: %v", err)
			}
			return fmt.Errorf("failed to generate voucher record: %v: %w", err, asynq.SkipRetry)
		}
		return nil
	} else {
		fmt.Println("campaign expired campExtId:", campU.CampaignExtID, "userExtId:", campU.UserExtID)
	}

	return nil
}
