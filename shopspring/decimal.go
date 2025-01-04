package shopspring

import (
	"fmt"
	"math"

	shopspring "github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc/basecalc"
)

var (
	ErrNewFromIntTooLarge = fmt.Errorf("creating decimal from integer: max allowed is %d", math.MaxInt64)
	ErrPowIntTooLarge     = fmt.Errorf("integer exponentiation: max allowed is %d", math.MaxInt32)
	ErrRoundTooLarge      = fmt.Errorf("rounding places:  max allowed is %d", math.MaxInt32)
)

type decimal struct {
	d            shopspring.Decimal
	divPrecision int32
}

var _ basecalc.Operator[decimal] = decimal{}

func newFromIntFunc(divPrecision int32) func(uint64) (decimal, error) {
	return func(n uint64) (decimal, error) {
		if n > math.MaxInt64 {
			return decimal{}, fmt.Errorf("%w: %d is too large", ErrNewFromIntTooLarge, n)
		}

		return decimal{
			d:            shopspring.New(int64(n), 0),
			divPrecision: divPrecision,
		}, nil
	}
}

func (d decimal) Mul(n decimal) (decimal, error) {
	return decimal{
		d:            d.d.Mul(n.d),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) Div(n decimal) (decimal, error) {
	return decimal{
		d:            d.d.DivRound(n.d, d.divPrecision),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) Sub(n decimal) (decimal, error) {
	return decimal{
		d:            d.d.Sub(n.d),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) Add(n decimal) (decimal, error) {
	return decimal{
		d:            d.d.Add(n.d),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) Abs() (decimal, error) {
	return decimal{
		d:            d.d.Abs(),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) LessThan(n decimal) (bool, error) {
	return d.d.LessThan(n.d), nil
}

func (d decimal) PowInt(n uint64) (decimal, error) {
	if n > math.MaxInt32 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrPowIntTooLarge, n)
	}

	res, err := d.d.PowInt32(int32(n))

	return decimal{
		d:            res,
		divPrecision: d.divPrecision,
	}, err
}

func (d decimal) Round(n uint64) (decimal, error) {
	if n > math.MaxInt32 {
		return decimal{}, fmt.Errorf("%w: %d is too large", ErrRoundTooLarge, n)
	}

	return decimal{
		d:            d.d.Round(int32(n)),
		divPrecision: d.divPrecision,
	}, nil
}

func (d decimal) IsNegative() (bool, error) {
	return d.d.IsNegative(), nil
}

func (d decimal) String() string {
	return d.d.String()
}
