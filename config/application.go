package config

type Application struct {
	Version     string `env:"VERSION"`
	Environment string `env:"ENVIRONMENT"`
}
