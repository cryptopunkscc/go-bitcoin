package lnd

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd/lnrpc"
)

var _ bitcoin.Lightning = &lightning{}

type lightning struct {
	rpc lnrpc.LightningClient
}

func (l lightning) Issue(ctx context.Context, invoice bitcoin.InvoiceRequest) (*bitcoin.Invoice, error) {
	res, err := l.rpc.AddInvoice(ctx, &lnrpc.Invoice{
		Value:  invoice.Amount.Sat(),
		Expiry: int64(invoice.ValidFor.Seconds()),
		Memo:   invoice.Memo,
	})
	if err != nil {
		return nil, err
	}
	li, err := l.rpc.LookupInvoice(ctx, &lnrpc.PaymentHash{
		RHash: res.RHash,
	})
	if err != nil {
		return nil, err
	}
	return lndInvoiceToInvoice(li), nil
}

func (l lightning) Balance(ctx context.Context) (bitcoin.Amount, error) {
	res, err := l.rpc.ChannelBalance(ctx, &lnrpc.ChannelBalanceRequest{})
	if err != nil {
		return 0, err
	}
	return bitcoin.Sat(res.Balance), nil
}

func (l lightning) Pay(ctx context.Context, invoice string) error {
	res, err := l.rpc.SendPaymentSync(ctx, &lnrpc.SendRequest{
		PaymentRequest: invoice,
	})
	if err != nil {
		return err
	}
	if res.PaymentError != "" {
		return errors.New(res.PaymentError)
	}
	return nil
}

func (l lightning) Decode(ctx context.Context, invoice string) *bitcoin.Invoice {
	res, err := l.rpc.DecodePayReq(ctx, &lnrpc.PayReqString{PayReq: invoice})
	if err != nil {
		return nil
	}
	return lndPayReqToInvoice(res)
}
