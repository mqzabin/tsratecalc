package shopspring

import (
	"errors"

	shopspring "github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc"
)

var (
	ErrConfigPrecisionNegative = errors.New("result precision must be positive")
	ErrRootNegative            = errors.New("root must be positive")
	ErrMaxTermsCacheNegative   = errors.New("max terms cache must be positive")
)

type Config struct {
	// Root is the root value to be used in the Taylor Series expansion.
	// It defines the "n" in the formula: "(1+x)^(1/n)-1".
	Root int32
	// Precision is the number of decimal places to consider in the calculations.
	// The calculation will only stop when the error is lower than 10^(-precision)/2.
	Precision int32
	// ConvergenceRadius sets the desired convergence radius for the rate value,
	// and will dynamically define how many Taylor Series terms will be used and pre-computed.
	//
	// The calculator will expand Taylor Series around x=0, until the convergence radius
	// boundaries (i.e. 0 + radius and 0 - radius) have error lower than the provided precision.
	//
	// It's recommended to be lower than 1. If it's greater than 1, the number of iterations (and Taylor terms cache)
	// required to converge on boundaries will grow exponentially.
	ConvergenceRadius shopspring.Decimal
	// MaxTermsCache is the maximum number of Taylor terms that will be cached.
	// If not provided, DefaultMaxTermsCache will be used.
	MaxTermsCache int32
}

// Calculator is a wrapper around tsratecalc.Calculator for "github.com/shopspring/decimal".Decimal type.
type Calculator struct {
	calc *tsratecalc.Calculator[decimal]
}

// NewCalculator creates a new calculator with the given Config.
func NewCalculator(cfg Config) (*Calculator, error) {
	if cfg.Precision < 0 {
		return nil, ErrConfigPrecisionNegative
	}

	if cfg.Root < 0 {
		return nil, ErrRootNegative
	}

	if cfg.MaxTermsCache < 0 {
		return nil, ErrMaxTermsCacheNegative
	}

	underlyingCfg := tsratecalc.Config[decimal]{
		Root:       uint64(cfg.Root),
		Precision:  uint64(cfg.Precision),
		NewFromInt: newFromIntFunc,
		ConvergenceRadius: decimal{
			cfg.ConvergenceRadius,
		},
		MaxTermsCache: uint64(cfg.MaxTermsCache),
	}

	calc, err := tsratecalc.NewCalculator[decimal](underlyingCfg)
	if err != nil {
		return nil, err
	}

	return &Calculator{
		calc: calc,
	}, nil
}

// ComputeRate receives a rate value and returns "(1+rate)^(1/root) - 1" using a Taylor Series expansion around rate=0.
// The root is defined in the calculator Config.
//
// The rate value should fall within the Config.ConvergenceRadius interval, around rate=0,
// otherwise ErrRateOutsideConvergenceBoundaries will be returned.
//
// It will return ConvergenceError if the desired precision is not achieved after the maximum number of iterations.
func (c *Calculator) ComputeRate(rate shopspring.Decimal) (shopspring.Decimal, error) {
	d := decimal{d: rate}

	result, err := c.calc.ComputeRate(d)
	if err != nil {
		return shopspring.Decimal{}, err
	}

	return result.d, nil
}

// TermsCacheLen returns the number of Taylor terms stored in the calculator's cache.
func (c *Calculator) TermsCacheLen() int {
	return c.calc.TermsCacheLen()
}
