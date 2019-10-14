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
	res, err := b.rpc.WalletBalance(ctx, &lnrpc.WalletBalanceRequest{})
	if err != nil {
		return 0, err
	}
	return bitcoin.Sat(res.ConfirmedBalance), nil
}

// NewAddress returns a new address for on-chain transactions
func (b blockchain) NewAddress(ctx context.Context, addrType int) (string, error) {
	t := lnrpc.AddressType_WITNESS_PUBKEY_HASH
	if addrType == bitcoin.NP2WKH {
		t = lnrpc.AddressType_NESTED_PUBKEY_HASH
	}
	res, err := b.rpc.NewAddress(ctx, &lnrpc.NewAddressRequest{
		Type: t,
	})
	if err != nil {
		return "", err
	}
	return res.Address, nil
}

// Send sends an amount of satoshis to the address with a given fee rate (0 for default). Returns a transaction ID.
func (b blockchain) Send(ctx context.Context, address string, amount bitcoin.Amount, feeRate int) (string, error) {
	res, err := b.rpc.SendCoins(ctx, &lnrpc.SendCoinsRequest{
		Addr:       address,
		Amount:     amount.Sat(),
		SatPerByte: int64(feeRate),
	})
	if err != nil {
		return "", err
	}
	return res.Txid, nil
}
