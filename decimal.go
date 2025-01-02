package tsratecalc

type Decimal interface {
	Mul(n Decimal) (Decimal, error)
	Div(n Decimal) (Decimal, error)
	Sub(n Decimal) (Decimal, error)
	Add(n Decimal) (Decimal, error)
	Abs() (Decimal, error)
	LessThan(n Decimal) (bool, error)
	PowInt(n uint32) (Decimal, error)
	Round(n uint32) (Decimal, error)
	IsNegative() (bool, error)
	String() string
}
