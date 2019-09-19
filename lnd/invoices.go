package lnd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnrpc"
)

func (lnd *Client) SetInvoiceHandler(h bitcoin.InvoiceHandler) {
	lnd.invoiceHandler = h
}

func (lnd *Client) CreateInvoice(amount bitcoin.Amount, memo string, expiry time.Duration) *bitcoin.Invoice {
	res, err := lnd.rpc.AddInvoice(context.Background(), &lnrpc.Invoice{
		Value:  amount.Sat(),
		Expiry: int64(expiry.Seconds()),
		Memo:   memo,
	})
	if err != nil {
		return nil
	}
	li, err := lnd.rpc.LookupInvoice(context.Background(), &lnrpc.PaymentHash{
		RHash: res.RHash,
	})
	if err != nil {
		return nil
	}
	return linvoiceToBinvoice(li)
}

func (lnd *Client) PayInvoice(invoice string) error {
	fmt.Println("Paying", invoice)
	res, err := lnd.rpc.SendPaymentSync(context.Background(), &lnrpc.SendRequest{
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

func (lnd *Client) subscribeInvoices() {
	s, err := lnd.rpc.SubscribeInvoices(context.Background(), &lnrpc.InvoiceSubscription{})
	if err != nil {
		panic(err)
	}
	for {
		update, err := s.Recv()
		if err != nil {
			return
		}
		if lnd.invoiceHandler != nil {
			lnd.invoiceHandler(linvoiceToBinvoice(update))
		}
	}
}

func (lnd *Client) ListInvoices() []*bitcoin.Invoice {
	res, err := lnd.rpc.ListInvoices(context.Background(), &lnrpc.ListInvoiceRequest{})
	if err != nil {
		return nil
	}
	list := make([]*bitcoin.Invoice, 0)
	for _, i := range res.GetInvoices() {
		list = append(list, linvoiceToBinvoice(i))
	}
	return list
}

func linvoiceToBinvoice(i *lnrpc.Invoice) *bitcoin.Invoice {
	return &bitcoin.Invoice{
		ID:             hex.EncodeToString(i.GetRPreimage()),
		IssuedAt:       time.Unix(i.GetCreationDate(), 0),
		ExpiresAt:      time.Unix(i.GetCreationDate()+i.GetExpiry(), 0),
		PaidAt:         time.Unix(i.GetSettleDate(), 0),
		AmountPaid:     bitcoin.Msat(i.GetAmtPaidMsat()),
		Amount:         bitcoin.Sat(i.GetValue()),
		Description:    i.GetMemo(),
		PaymentRequest: i.GetPaymentRequest(),
		State:          int(i.State),
	}
}

func (lnd *Client) DecodeInvoice(req string) *bitcoin.Invoice {
	res, err := lnd.rpc.DecodePayReq(context.Background(), &lnrpc.PayReqString{PayReq: req})
	if err != nil {
		return nil
	}
	return payreqToBinvoice(res)
}

func payreqToBinvoice(req *lnrpc.PayReq) *bitcoin.Invoice {
	return &bitcoin.Invoice{
		Description: req.GetDescription(),
		IssuedAt:    time.Unix(req.GetTimestamp(), 0),
		ExpiresAt:   time.Unix(req.GetTimestamp()+req.GetExpiry(), 0),
		Amount:      bitcoin.Sat(req.GetNumSatoshis()),
	}
}
