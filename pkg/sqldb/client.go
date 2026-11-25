package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient opens a Postgres connection using the settings read from env vars with the given prefix.
func NewClient(envPrefix string) (*gorm.DB, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return db, nil
}
