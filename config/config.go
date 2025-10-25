package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
)

var (
	instance *Config
	once     sync.Once
	iniErr   error
)

type Config struct {
	Server      Server
	Postgresql  Postgresql
	Application Application
	RateLimiter RateLimiter
}

func NewConfig() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		iniErr = env.Parse(instance)
		if iniErr != nil {
			instance = nil
		}
	})
	return instance, iniErr
}
