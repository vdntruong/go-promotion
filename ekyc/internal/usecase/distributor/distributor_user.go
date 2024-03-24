package distributor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const EventTypeUserFirstLogin = "user:first-time-login"

type UserFirstLoginPayload struct {
	UserExtID     string
	LoginDateTime time.Time
}

func newUserFirstLoginEvent(userExtID string, dateTime time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(UserFirstLoginPayload{UserExtID: userExtID, LoginDateTime: dateTime})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(EventTypeUserFirstLogin, payload), nil
}

type UserDistributor struct {
	client *asynq.Client
}

func NewUserDistributor(c *asynq.Client) *UserDistributor {
	return &UserDistributor{
		client: c,
	}
}
func (u *UserDistributor) DispatchUserFirstTimeLogin(userExtID string, dateTime time.Time) error {
	t, err := newUserFirstLoginEvent(userExtID, dateTime)
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
