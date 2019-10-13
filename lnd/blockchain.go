package lnd

import (
	"context"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd/lnrpc"
)

var _ bitcoin.Blockchain = &blockchain{}

type blockchain struct {
	rpc lnrpc.LightningClient
}

// Balance returns confirmed (spendable) on-chain balance
func (b blockchain) Balance(ctx context.Context) (bitcoin.Amount, error) {
	res, err := b.rpc.WalletBalance(ctx, &lnrpc.WalletBalanceRequest{}, nil)
	if err != nil {
		return 0, err
	}
	return bitcoin.Sat(res.ConfirmedBalance), nil
}
