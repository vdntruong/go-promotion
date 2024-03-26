package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"ekyc/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        UserRepo
	auth        Authenticator
	distributor UserEventDistributor
}

func NewUserService(r UserRepo, a Authenticator, d UserEventDistributor) *UserService {
	return &UserService{
		repo:        r,
		auth:        a,
		distributor: d,
	}
}

func (c *UserService) SignUp(ctx context.Context, email, password string, campaignExtID *string) (*model.User, error) {
	_, found, err := c.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if found {
		return nil, model.NewRespError(
			errors.New("user existing"),
			model.WithStatusCode(http.StatusConflict),
			model.WithDescription("Email is in use"),
		)
	}

	hashedPass, err := generatePassword(password)
	if err != nil {
		return nil, err
	}

	var u = model.User{
		Email:         email,
		Password:      hashedPass,
		CampaignExtID: campaignExtID,
	}
	result, err := c.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	u, found, err := c.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !found {
		return "", model.NewRespError(
			errors.New("email is not exist"),
			model.WithStatusCode(http.StatusNotFound),
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", model.NewRespError(
			errors.New("invalid user password"),
			model.WithStatusCode(http.StatusBadRequest),
		)
	}

	if u.FirstLoginDate == nil {
		var firstLoginTime = time.Now()
		u.FirstLoginDate = &firstLoginTime
		var callback = func() error { return nil }

		// event for user registered based on a campaign and just have fist login
		if u.CampaignExtID != nil {
			callback = func() error {
				return c.distributor.DispatchUserFirstTimeLogin(u.ExtID, *u.CampaignExtID, u.CreatedAt, firstLoginTime)
			}
		}

		if err := c.repo.UpdateFirstLogin(ctx, u.ID, firstLoginTime, callback); err != nil {
			return "", err
		}
	}
	return c.auth.GetToken(*u)
}

func generatePassword(pass string) (string, error) {
	rs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}
	return string(rs), nil
}
