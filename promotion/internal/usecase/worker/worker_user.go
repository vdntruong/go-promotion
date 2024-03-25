package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"promotion/internal/model"
	"promotion/internal/usecase"

	"github.com/hibiken/asynq"
)

type UserTaskHandler struct {
	uc     usecase.CampaignUserUsecase
	campUc usecase.CampaignUsecase
}

func NewUserTaskHandler(c usecase.CampaignUserUsecase, campUc usecase.CampaignUsecase) *UserTaskHandler {
	return &UserTaskHandler{uc: c, campUc: campUc}
}

func (u *UserTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p UserFirstLoginPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return u.handleUserFirstTimeLogin(ctx, p)
}

func (u *UserTaskHandler) handleUserFirstTimeLogin(ctx context.Context, p UserFirstLoginPayload) error {
	campU, camp, err := u.uc.CreateCampaignUser(ctx, p.CampaignExtID, p.UserExtID, p.LoginDateTime)
	if err != nil {
		if errors.Is(err, model.ErrDBNotFound) || errors.Is(err, model.ErrDuplicate) {
			return fmt.Errorf("failed to add campaign user record: %v: %w", err, asynq.SkipRetry)
		}
		return fmt.Errorf("failed to add campaign user record: %v", err)
	}

	if campU.FirstLoginDate.After(camp.StartDate) && campU.FirstLoginDate.Before(camp.EndDate) {
		fmt.Println("call to voucher to generate voucher campid:", campU.CampaignExtID, "user:", campU.UserExtID)
		return nil
	} else {
		fmt.Println("generate voucher expired campid:", campU.CampaignExtID, "user:", campU.UserExtID)
	}

	return nil
}
