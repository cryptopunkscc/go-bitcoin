package lnd

import (
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Host        string
	Port        int
	Macaroon    []byte
	Certificate []byte
}

func (cfg *Config) validate() error {
	if cfg.Host == "" {
		return errors.New("host missing")
	}
	if (cfg.Port < 1) || (cfg.Port > 65535) {
		return errors.New("port invalid")
	}
	if len(cfg.Macaroon) == 0 {
		return errors.New("macaroon missing")
	}
	if len(cfg.Certificate) == 0 {
		return errors.New("cert missing")
	}
	return nil
}

func (cfg *Config) macaroonCreds() MacaroonCredential {
	return MacaroonCredential{cfg.Macaroon}
}

func (cfg *Config) tlsCreds() credentials.TransportCredentials {
	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(cfg.Certificate) {
		return nil
	}

	return credentials.NewClientTLSFromCert(certPool, "")
}

func (cfg *Config) url() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
