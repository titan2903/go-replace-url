package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	conf *Config
)

type Config struct {
	Port                 int           `envconfig:"PORT" default:"8080"`
	DB                   string        `envconfig:"DB_NAME" default:"emi_site"`
	DBHost               string        `envconfig:"DB_HOST" default:"localhost"`
	DBPort               int           `envconfig:"DB_PORT" default:"5432"`
	DBUsername           string        `envconfig:"DB_USERNAME" default:"postgres"`
	DBPassword           string        `envconfig:"DB_PASSWORD" default:"password"`
	DBConnectionIdle     time.Duration `envconfig:"DB_CONNECTION_IDLE" default:"1m"`
	DBConnectionLifetime time.Duration `envconfig:"DB_CONNECTION_LIFETIME" default:"5m"`
	DBMaxIdle            int           `envconfig:"DB_MAX_IDLE" default:"30"`
	DBMaxOpen            int           `envconfig:"DB_MAX_OPEN" default:"90"`
}

func Init() {
	conf = new(Config)
	envconfig.MustProcess("", conf)
}

func Get() *Config {
	return conf
}
