package shopspring

import (
	"errors"

	shopspring "github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc/basecalc"
)

var (
	ErrConfigPrecisionNegative = errors.New("result precision must be positive")
	ErrRootNegative            = errors.New("root must be positive")
	ErrMaxIterationsNegative   = errors.New("max iterations must be positive")
)

type Config struct {
	Root            int32
	MaxIterations   int32
	ResultPrecision int32
	DivPrecision    int32
}

type Calculator struct {
	calc         *basecalc.Calculator[decimal]
	divPrecision int32
}

func NewCalculator(cfg Config) (*Calculator, error) {
	if cfg.ResultPrecision < 0 {
		return nil, ErrConfigPrecisionNegative
	}

	if cfg.Root < 0 {
		return nil, ErrRootNegative
	}

	if cfg.MaxIterations < 0 {
		return nil, ErrMaxIterationsNegative
	}

	underlyingCfg := basecalc.Config[decimal]{
		Root:          uint64(cfg.Root),
		Precision:     uint64(cfg.ResultPrecision),
		NewFromInt:    newFromIntFunc(cfg.DivPrecision),
		MaxIterations: uint64(cfg.MaxIterations),
	}

	calc, err := basecalc.NewCalculator[decimal](underlyingCfg)
	if err != nil {
		return nil, err
	}

	return &Calculator{
		calc:         calc,
		divPrecision: cfg.DivPrecision,
	}, nil
}

func (c *Calculator) ComputeRate(rate shopspring.Decimal) (shopspring.Decimal, error) {
	d := decimal{d: rate, divPrecision: c.divPrecision}

	result, err := c.calc.ComputeRate(d)
	if err != nil {
		return shopspring.Decimal{}, err
	}

	return result.d, nil
}
