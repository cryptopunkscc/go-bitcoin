package bitcoin

import "context"

type Wallet interface {
	Network(context.Context) (string, error)
	Agent(context.Context) (string, error)
	Blockchain() Blockchain
	Lightning() Lightning
}

type Blockchain interface {
	Balance(context.Context) (Amount, error)
}

type Lightning interface {
	Balance(context.Context) (Amount, error)
	Pay(context.Context, string) error
	Decode(context.Context, string) *Invoice
	Issue(context.Context, InvoiceRequest) (*Invoice, error)
}
