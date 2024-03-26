package worker

import "time"

const EventTypeUserFirstLogin = "user:campaign:first-time-login"

type UserFirstLoginPayload struct {
	UserExtID     string
	CampaignExtID string
	RegisterTime  time.Time
	LoginDateTime time.Time
}
