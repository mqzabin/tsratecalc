package tsratecalc

type Operator[Decimal any] interface {
	Mul(n Decimal) (Decimal, error)
	DivRound(n Decimal, places uint64) (Decimal, error)
	Sub(n Decimal) (Decimal, error)
	Add(n Decimal) (Decimal, error)
	Abs() (Decimal, error)
	LessThanOrEqual(n Decimal) (bool, error)
	PowInt(n uint64) (Decimal, error)
	Truncate(places uint64) (Decimal, error)
	String() string
}
