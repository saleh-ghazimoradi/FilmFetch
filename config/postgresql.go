package config

import "time"

type Postgresql struct {
	Host        string        `env:"POSTGRES_HOST"`
	Port        string        `env:"POSTGRES_PORT"`
	User        string        `env:"POSTGRES_USER"`
	Password    string        `env:"POSTGRES_PASSWORD"`
	Name        string        `env:"POSTGRES_NAME"`
	MaxOpenConn int           `env:"POSTGRES_MAX_OPEN_CONN"`
	MaxIdleConn int           `env:"POSTGRES_MAX_IDLE_CONN"`
	MaxIdleTime time.Duration `env:"POSTGRES_MAX_IDLE_TIME"`
	SSLMode     string        `env:"POSTGRES_SSL_MODE"`
	Timeout     time.Duration `env:"POSTGRES_TIMEOUT"`
}
