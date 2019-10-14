package bitcoin

import "context"

type Wallet interface {
	Network(context.Context) (string, error)
	Agent(context.Context) (string, error)
	Blockchain() Blockchain
	Lightning() Lightning
}

const (
	DEFAULT = iota
	P2WKH
	NP2WKH
)

type Blockchain interface {
	Balance(context.Context) (Amount, error)
	NewAddress(ctx context.Context, addrType int) (string, error)
	Send(ctx context.Context, address string, amount Amount, feeRate int) (string, error)
}

type Lightning interface {
	Balance(context.Context) (Amount, error)
	Pay(context.Context, string) error
	Decode(context.Context, string) *Invoice
	Issue(context.Context, InvoiceRequest) (*Invoice, error)
}
