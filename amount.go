package bitcoin

type Amount int64

const msatsPerBtc = 100_000_000_000

func Sat(sats int64) Amount {
	return Amount(sats * 1000)
}

func Msat(msats int64) Amount {
	return Amount(msats)
}

func Btc(btc float64) Amount {
	return Amount(btc * msatsPerBtc)
}

func (a Amount) Sat() int64 {
	return int64(a / 1000)
}

func (a Amount) Msat() int64 {
	return int64(a)
}

func (a Amount) Btc() float64 {
	return float64(a) / msatsPerBtc
}
