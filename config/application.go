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
