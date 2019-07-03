package bitcoin

import "time"

type LightningClient interface {
	// Generic stuff
	Version() string
	Balance() int
	SetChannelHandler(ChannelHandler)

	// Invoicing
	CreateInvoice(amount Amount, memo string, expiry time.Duration) *Invoice
	PayInvoice(invoice string) error
	ListInvoices() []*Invoice
	DecodeInvoice(string) *Invoice
	SetInvoiceHandler(InvoiceHandler)
}
