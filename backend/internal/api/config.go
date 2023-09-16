package api

import (
	"fmt"

	"github.com/Netflix/go-env"
)

type ServerConfig struct {
	Port uint16 `env:"DASHBOARD_API_PORT,default=3000"`
}

func LoadServerConfig() (*ServerConfig, error) {
	config := &ServerConfig{}
	_, err := env.UnmarshalFromEnviron(config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}
	return config, nil
}

func (c *ServerConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("port cannot be <= 0 (got: %d)", c.Port)
	}
	return nil
}
