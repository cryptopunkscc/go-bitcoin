package bitcoin

import "time"

type Invoice struct {
	ID             string
	Description    string
	Amount         Amount
	AmountPaid     Amount
	IssuedAt       time.Time
	ExpiresAt      time.Time
	PaidAt         time.Time
	PaymentRequest string
	State          int
}

type InvoiceHandler func(*Invoice)

func (i Invoice) Paid() bool {
	return i.PaidAt != time.Time{}
}
