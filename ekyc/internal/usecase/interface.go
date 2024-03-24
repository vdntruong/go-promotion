package usecase

import (
	"context"
	"time"

	"ekyc/internal/model"
)

type Authenticator interface {
	GetToken(model.User) (string, error)
	VerifyToken(string) error
}

type UserUsecase interface {
	SignUp(ctx context.Context, email, password string, campaignExtID string) (*model.User, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type UserEventDistributor interface {
	DispatchUserFirstTimeLogin(userExtID string, dateTime time.Time) error
}

type UserRepo interface {
	CreateUser(context.Context, model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, bool, error)
	UpdateFirstLogin(ctx context.Context, id uint, dateTime time.Time, callback func() error) error
}
