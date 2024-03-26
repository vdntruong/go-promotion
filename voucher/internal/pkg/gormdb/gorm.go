package gormdb

import (
	"fmt"

	"voucher/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}