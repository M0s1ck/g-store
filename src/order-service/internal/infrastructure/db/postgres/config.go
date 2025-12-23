package postgres

import (
	"time"

	"orders-service/internal/config"
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
		Host:     cfg.OrdersDB.Host,
		Port:     cfg.OrdersDB.Port,
		Name:     cfg.OrdersDB.Name,
		User:     cfg.OrdersDB.User,
		Password: cfg.OrdersDB.Password,
		SSLMode:  cfg.OrdersDB.SSLMode,

		MaxOpenConns:    cfg.OrdersDB.MaxOpenConns,
		MaxIdleConns:    cfg.OrdersDB.MaxIdleConns,
		ConnMaxLifetime: cfg.OrdersDB.ConnMaxLifetime,
	}
}
