package shopspring

import (
	"fmt"
	"math"

	shopspring "github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc"
)

var (
	ErrNewFromIntTooLarge = fmt.Errorf("creating decimal from integer: max allowed is %d", math.MaxInt64)
	ErrPowIntTooLarge     = fmt.Errorf("integer exponentiation: max allowed is %d", math.MaxInt32)
	ErrRoundTooLarge      = fmt.Errorf("rounding places:  max allowed is %d", math.MaxInt32)
)

type decimal struct {
	d shopspring.Decimal
}

var _ tsratecalc.Operator[decimal] = decimal{}

func newFromIntFunc(n uint64) (decimal, error) {

	if n > math.MaxInt64 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrNewFromIntTooLarge, n)
	}

	return decimal{
		d: shopspring.New(int64(n), 0),
	}, nil
}

func (d decimal) Mul(n decimal) (decimal, error) {
	return decimal{
		d: d.d.Mul(n.d),
	}, nil
}

func (d decimal) DivRound(n decimal, places uint64) (decimal, error) {
	if places > math.MaxInt32 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrRoundTooLarge, places)
	}

	return decimal{
		d: d.d.DivRound(n.d, int32(places)),
	}, nil
}

func (d decimal) Sub(n decimal) (decimal, error) {
	return decimal{
		d: d.d.Sub(n.d),
	}, nil
}

func (d decimal) Add(n decimal) (decimal, error) {
	return decimal{
		d: d.d.Add(n.d),
	}, nil
}

func (d decimal) Abs() (decimal, error) {
	return decimal{
		d: d.d.Abs(),
	}, nil
}

func (d decimal) LessThanOrEqual(n decimal) (bool, error) {
	return d.d.LessThanOrEqual(n.d), nil
}

func (d decimal) PowInt(n uint64) (decimal, error) {
	if n > math.MaxInt32 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrPowIntTooLarge, n)
	}

	res, err := d.d.PowInt32(int32(n))

	return decimal{
		d: res,
	}, err
}

func (d decimal) Truncate(n uint64) (decimal, error) {
	if n > math.MaxInt32 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrRoundTooLarge, n)
	}

	return decimal{
		d: d.d.Truncate(int32(n)),
	}, nil
}

func (d decimal) IsNegative() (bool, error) {
	return d.d.IsNegative(), nil
}

func (d decimal) String() string {
	return d.d.String()
}
