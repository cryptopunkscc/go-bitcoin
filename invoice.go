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

type InvoiceRequest struct {
	Amount   Amount
	Memo     string
	ValidFor time.Duration
}

func (i Invoice) Paid() bool {
	return i.PaidAt != time.Time{}
}

func (i Invoice) ExpiresIn() time.Duration {
	now := time.Now()
	if now.After(i.ExpiresAt) {
		return 0
	}
	return i.ExpiresAt.Sub(now)
}
