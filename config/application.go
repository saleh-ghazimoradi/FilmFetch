package config

type Application struct {
	Version     string `env:"VERSION"`
	Environment string `env:"ENVIRONMENT"`
}

type RateLimiter struct {
	RPS     float64 `env:"RPS"`
	Burst   int     `env:"BURST"`
	Enabled bool    `env:"ENABLED"`
}

type Mail struct {
	Host     string `env:"MAIL_HOST"`
	Port     int    `env:"MAIL_PORT"`
	Username string `env:"MAIL_USERNAME"`
	Password string `env:"MAIL_PASSWORD"`
	Sender   string `env:"MAIL_SENDER"`
}
