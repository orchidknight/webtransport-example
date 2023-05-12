package config

import (
	"fmt"
	"github.com/caarlos0/env"
)

type Config struct {
	Port                string `env:"HTTP_PORT"`
	Host                string `env:"HTTP_HOST"`
	CertificatePath     string `env:"CERTIFICATE_PATH"`
	KeyPath             string `env:"KEY_PATH"`
	PingDelaySeconds    int    `env:"PING_DELAY_SECONDS"`
	WebTransportLogFile string `env:"WEB_TRANSPORT_LOG_FILE"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("wrong port")
	}

	return cfg, nil
}
