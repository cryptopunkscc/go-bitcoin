package lnd

import (
	"context"
	"errors"
	"google.golang.org/grpc"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd/lnrpc"
)

// Verify that Wallet satisfies bitcoin.Wallet interface
var _ bitcoin.Wallet = &Wallet{}

type Wallet struct {
	rpc        lnrpc.LightningClient
	blockchain *blockchain
	lightning  *lightning
}

// New returns a new instance of an LND wallet
func New(cfg *Config) (*Wallet, error) {
	// check basic validity of config
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// make a grpc connection
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cfg.tlsCreds()),
		grpc.WithPerRPCCredentials(cfg.macaroonCreds()),
	}
	conn, err := grpc.Dial(cfg.url(), opts...)
	if err != nil {
		return nil, err
	}

	rpc := lnrpc.NewLightningClient(conn)
	wallet := &Wallet{
		rpc:        rpc,
		blockchain: &blockchain{rpc: rpc},
		lightning:  &lightning{rpc: rpc},
	}
	return wallet, nil
}

// Network returns the network agent is connected to (mainnet, testnet, regtest)
func (wallet *Wallet) Network(ctx context.Context) (string, error) {
	info, err := wallet.rpc.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return "", err
	}
	for _, chain := range info.Chains {
		if chain.Chain == "bitcoin" {
			return chain.Network, nil
		}
	}
	return "", errors.New("unknown chain")
}

// Agent returns the agent name and version
func (wallet *Wallet) Agent(ctx context.Context) (string, error) {
	info, err := wallet.rpc.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return "", err
	}
	return "wallet " + info.Version, nil
}

// Blockchain returns an onchain client if available
func (wallet *Wallet) Blockchain() bitcoin.Blockchain {
	return wallet.blockchain
}

// Lightning returns a lightning client if available
func (wallet *Wallet) Lightning() bitcoin.Lightning {
	return wallet.lightning
}
