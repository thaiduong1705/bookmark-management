package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetDSN builds the Postgres connection string.
func (c *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", c.Host, c.User, c.Password, c.DBName, c.Port)
}
