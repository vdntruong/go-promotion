package distributor

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const EventTypeUserFirstLogin = "user:campaign:first-time-login"

type UserFirstLoginPayload struct {
	UserExtID     string
	CampaignExtID string
	RegisterTime  time.Time
	LoginDateTime time.Time
}

func newUserFirstLoginEvent(userExtID string, campaignExtID string, register, login time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(UserFirstLoginPayload{
		UserExtID:     userExtID,
		CampaignExtID: campaignExtID,
		RegisterTime:  register,
		LoginDateTime: login,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(EventTypeUserFirstLogin, payload), nil
}
