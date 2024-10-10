package postgres

import (
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/config"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/db"
	"gorm.io/gorm"
)

func NewFromConfig(cfg config.KVStore, path string) (*gorm.DB, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := cfg.GetInt("MAX_IDLE_CONNECTION"); maxIdle > 0 {
		connOpt.MaxIdle = cfg.GetInt("MAX_IDLE_CONNECTION")
	}

	if maxOpen := cfg.GetInt("MAX_OPEN_CONNECTION"); maxOpen > 0 {
		connOpt.MaxOpen = maxOpen
	}

	if maxLifetime := cfg.GetDuration("MAX_LIFETIME_CONNECTION"); maxLifetime > 0 {
		connOpt.MaxLifetime = maxLifetime
	}

	if connTimeout := cfg.GetDuration("TIMEOUT"); connTimeout > 0 {
		connOpt.ConnectTimeout = connTimeout
	}

	if keepAlive := cfg.GetDuration("KEEP_ALIVE_INTERVAL"); keepAlive > 0 {
		connOpt.KeepAliveCheckInterval = keepAlive
	}

	opt, err := db.NewDatabaseOption(
		cfg.GetString("DB_HOST"),
		cfg.GetInt("DB_PORT"),
		cfg.GetString("DB_USERNAME"),
		cfg.GetString("DB_PASSWORD"),
		cfg.GetString("DB_NAME"),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	return NewPostgresDatabase(opt)
}
