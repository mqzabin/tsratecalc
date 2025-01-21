package tsratecalc

// Operator defines the interface to an arbitrary/fixed precision decimal type to be used by the Calculator,
// and compute (1+x)^(1/n)-1 using Taylor series expansion around x=0.
//
// Most operations could return an error if the operation is not possible for some reason (e.g. overflows on fixed precision decimals).
type Operator[Decimal any] interface {
	// Mul multiply two decimals.
	Mul(n Decimal) (Decimal, error)
	// DivRound divide two decimals and round the result to the provided number of decimal places.
	DivRound(n Decimal, places uint64) (Decimal, error)
	// Sub subtracts two decimals.
	Sub(n Decimal) (Decimal, error)
	// Add adds two decimals.
	Add(n Decimal) (Decimal, error)
	// Abs returns the absolute value of the decimal.
	Abs() (Decimal, error)
	// LessThanOrEqual returns true if the decimal is less than or equal to the provided decimal.
	LessThanOrEqual(n Decimal) (bool, error)
	// PowInt returns the decimal raised to the power of the provided integer.
	PowInt(n uint64) (Decimal, error)
	// Truncate returns the decimal truncated to the provided number of decimal places.
	Truncate(places uint64) (Decimal, error)
	// String returns the string representation of the decimal.
	String() string
}
