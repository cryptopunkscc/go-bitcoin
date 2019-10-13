package lnd

import (
	"encoding/hex"
	"fmt"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd/lnrpc"
	"strconv"
	"strings"
	"time"
)

func bolt7StringToUint64(id string) uint64 {
	split := strings.Split(id, ":")
	var res uint64

	if len(split) != 3 {
		return 0
	}

	height, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return 0
	}

	txidx, err := strconv.ParseUint(split[1], 10, 64)
	if err != nil {
		return 0
	}

	outputidx, err := strconv.ParseUint(split[2], 10, 64)
	if err != nil {
		return 0
	}

	res = ((height & 0xffffff) << 40) | ((txidx & 0xffffff) << 16) | (outputidx & 0xffff)

	return res
}

func bolt7Uint64ToString(id uint64) string {
	var height = uint64(id >> 40)
	var txidx = uint64((id & 0xffffff0000) >> 16)
	var outputidx = uint64(id & 0xffff)

	return fmt.Sprintf("%d:%d:%d", height, txidx, outputidx)
}

// txIDBytesToString converts byte representation of txid to string.
// Bytes need to be reversed before encoding to hex due to historical reasons.
// Details: http://learnmeabitcoin.com/glossary/txid#why
func txIDBytesToString(txid []byte) string {
	for i, j := 0, len(txid)-1; i < j; i, j = i+1, j-1 {
		txid[i], txid[j] = txid[j], txid[i]
	}

	return hex.EncodeToString(txid)
}

// txIDBytesToString converts a string representation of txid to bytes.
// Bytes need to be reversed after decoding from hex due to historical reasons.
// Details: http://learnmeabitcoin.com/glossary/txid#why
func txIDStringToBytes(txid string) []byte {
	bytes, err := hex.DecodeString(txid)

	if err != nil {
		return []byte{}
	}

	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return bytes
}

func lndInvoiceToInvoice(i *lnrpc.Invoice) *bitcoin.Invoice {
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

func lndPayReqToInvoice(req *lnrpc.PayReq) *bitcoin.Invoice {
	return &bitcoin.Invoice{
		Description: req.GetDescription(),
		IssuedAt:    time.Unix(req.GetTimestamp(), 0),
		ExpiresAt:   time.Unix(req.GetTimestamp()+req.GetExpiry(), 0),
		Amount:      bitcoin.Sat(req.GetNumSatoshis()),
	}
}
