package shopspring

import (
	"errors"

	shopspring "github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc"
)

var (
	ErrConfigPrecisionNegative = errors.New("result precision must be positive")
	ErrRootNegative            = errors.New("root must be positive")
)

type Config struct {
	Root              int32
	Precision         int32
	ConvergenceRadius shopspring.Decimal
}

type Calculator struct {
	calc *tsratecalc.Calculator[decimal]
}

func NewCalculator(cfg Config) (*Calculator, error) {
	if cfg.Precision < 0 {
		return nil, ErrConfigPrecisionNegative
	}

	if cfg.Root < 0 {
		return nil, ErrRootNegative
	}

	underlyingCfg := tsratecalc.Config[decimal]{
		Root:       uint64(cfg.Root),
		Precision:  uint64(cfg.Precision),
		NewFromInt: newFromIntFunc,
		ConvergenceRadius: decimal{
			cfg.ConvergenceRadius,
		},
	}

	calc, err := tsratecalc.NewCalculator[decimal](underlyingCfg)
	if err != nil {
		return nil, err
	}

	return &Calculator{
		calc: calc,
	}, nil
}

func (c *Calculator) ComputeRate(rate shopspring.Decimal) (shopspring.Decimal, error) {
	d := decimal{d: rate}

	result, err := c.calc.ComputeRate(d)
	if err != nil {
		return shopspring.Decimal{}, err
	}

	return result.d, nil
}

func (c *Calculator) TermsCacheLen() int {
	return c.calc.TermsCacheLen()
}
