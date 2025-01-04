package basecalc

type Operator[Decimal any] interface {
	Mul(n Decimal) (Decimal, error)
	Div(n Decimal) (Decimal, error)
	Sub(n Decimal) (Decimal, error)
	Add(n Decimal) (Decimal, error)
	Abs() (Decimal, error)
	LessThan(n Decimal) (bool, error)
	PowInt(n uint64) (Decimal, error)
	Round(places uint64) (Decimal, error)
	IsNegative() (bool, error)
	String() string
}
