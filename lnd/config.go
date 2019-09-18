package lnd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	MacaroonPath string `json:"macaroon"`
	CertPath     string `json:"cert"`
}

func (cfg *Config) Validate() error {
	if cfg.Host == "" {
		return errors.New("LND host missing")
	}
	if (cfg.Port < 1) || (cfg.Port > 65535) {
		return errors.New("LND port invalid")
	}
	if cfg.MacaroonPath == "" {
		return errors.New("LND macaroon path missing")
	}
	if cfg.CertPath == "" {
		return errors.New("LND cert path missing")
	}
	return nil
}

func (cfg *Config) Macaroon() MacaroonCredential {
	bytes, err := ioutil.ReadFile(cfg.MacaroonPath)

	if err != nil {
		return MacaroonCredential{}
	}

	return MacaroonCredential{bytes}
}

func (cfg *Config) TLSCredentials() credentials.TransportCredentials {
	tlsCreds, err := credentials.NewClientTLSFromFile(cfg.CertPath, "")

	if err != nil {
		return nil
	}

	return tlsCreds
}

func (cfg *Config) url() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
