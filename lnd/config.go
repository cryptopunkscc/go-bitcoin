package lnd

import (
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
