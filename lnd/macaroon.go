package lnd

import (
	"context"
	"encoding/hex"
)

// MacaroonCredential implements PerRPCCredentials interface
type MacaroonCredential struct {
	Bytes []byte
}

// RequireTransportSecurity implements PerRPCCredentials interface
func (m MacaroonCredential) RequireTransportSecurity() bool {
	return true
}

// GetRequestMetadata returnt the hex-encoded macaroon used as RPC credential
func (m MacaroonCredential) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {

	md := make(map[string]string)
	md["macaroon"] = hex.EncodeToString(m.Bytes)
	return md, nil
}
