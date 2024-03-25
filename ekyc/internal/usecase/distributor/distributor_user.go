package distributor

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type UserDistributor struct {
	client *asynq.Client
}

func NewUserDistributor(c *asynq.Client) *UserDistributor {
	return &UserDistributor{
		client: c,
	}
}

func (u *UserDistributor) DispatchUserFirstTimeLogin(userExtID string, campaignExtID string, dateTime time.Time) error {
	t, err := newUserFirstLoginEvent(userExtID, campaignExtID, dateTime)
	if err != nil {
		return err
	}

	if info, err := u.client.Enqueue(t); err != nil {
		return err
	} else {
		log.Printf("enqueued user first time login event: id=%s queue=%s", info.ID, info.Queue)
	}
	return nil
}
