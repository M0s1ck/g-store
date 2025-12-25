package postgres

import (
	"time"

	"payment-service/internal/config"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Host:     cfg.PaymentsDB.Host,
		Port:     cfg.PaymentsDB.Port,
		Name:     cfg.PaymentsDB.Name,
		User:     cfg.PaymentsDB.User,
		Password: cfg.PaymentsDB.Password,
		SSLMode:  cfg.PaymentsDB.SSLMode,

		MaxOpenConns:    cfg.PaymentsDB.MaxOpenConns,
		MaxIdleConns:    cfg.PaymentsDB.MaxIdleConns,
		ConnMaxLifetime: cfg.PaymentsDB.ConnMaxLifetime,
	}
}
