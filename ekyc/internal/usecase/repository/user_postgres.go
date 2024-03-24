package repository

import (
	"context"
	"errors"
	"time"

	"ekyc/internal/model"
	"ekyc/internal/pkg/util"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &UserRepository{db: db}
}

func (c *UserRepository) CreateUser(ctx context.Context, u model.User) (*model.User, error) {
	if err := c.db.WithContext(ctx).Create(&u).Error; err != nil {
		if util.IsDuplicateError(err) {
			return nil, model.NewRespError(err, model.WithDescription("Email is in use"))
		}
		return nil, err
	}
	return &u, nil
}

func (c *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, bool, error) {
	var result *model.User
	if err := c.db.WithContext(ctx).First(&result, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return result, true, nil
}

func (c *UserRepository) UpdateFirstLogin(ctx context.Context, id uint, dateTime time.Time, callback func() error) error {
	tx := c.db.Begin()
	
	var result = model.User{Model: gorm.Model{ID: id}, FirstLoginDate: &dateTime}
	if err := tx.WithContext(ctx).Model(&result).
		Where("id = ?", id).
		Update("first_login_date", &dateTime).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := callback(); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
