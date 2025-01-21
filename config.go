package tsratecalc

import (
	"errors"
	"fmt"
)

const (
	minPrecision = 0
	minRoot      = 2

	// DefaultMaxTermsCache defines the default maximum number of terms to cache.
	DefaultMaxTermsCache = 30000
)

var (
	ErrConfigRootMinValue              = fmt.Errorf("root should be greater than %d", minRoot)
	ErrConfigPrecisionMinValue         = fmt.Errorf("precision should be greater than %d", minPrecision)
	ErrConfigNewFromIntIsNil           = errors.New("'decimal from integer' factory should not be nil")
	ErrConfigConvergenceRadiusPositive = errors.New("convergence radius must be positive")
)

type Config[Decimal Operator[Decimal]] struct {
	// Root is the root value to be used in the Taylor Series expansion.
	// It defines the "n" in the formula: "(1+x)^(1/n)-1".
	Root uint64

	// Precision is the number of decimal places to consider in the calculations.
	// The calculation will only stop when the error is lower than 10^(-precision)/2.
	Precision uint64

	// NewFromInt is a factory function that creates a Decimal from an integer.
	NewFromInt func(n uint64) (Decimal, error)

	// ConvergenceRadius sets the desired convergence radius for the rate value,
	// and will dynamically define how many Taylor Series terms will be used and pre-computed.
	//
	// The calculator will expand Taylor Series around x=0, until the convergence radius
	// boundaries (i.e. 0 + radius and 0 - radius) have error lower than the provided precision.
	//
	// It's recommended to be lower than 1. If it's greater than 1, the number of iterations (and Taylor terms cache)
	// required to converge on boundaries will grow exponentially.
	ConvergenceRadius Decimal

	// MaxTermsCache is the maximum number of Taylor terms that will be cached.
	// If not provided, DefaultMaxTermsCache will be used.
	MaxTermsCache uint64
}

func validateConfig[Decimal Operator[Decimal]](cfg Config[Decimal]) (Config[Decimal], error) {
	if cfg.Root < minRoot {
		return Config[Decimal]{}, ErrConfigRootMinValue
	}

	if cfg.Precision < minPrecision {
		return Config[Decimal]{}, ErrConfigPrecisionMinValue
	}

	if cfg.NewFromInt == nil {
		return Config[Decimal]{}, ErrConfigNewFromIntIsNil
	}

	zero, err := cfg.NewFromInt(0)
	if err != nil {
		return Config[Decimal]{}, fmt.Errorf("creating '0' decimal: %w", err)
	}

	neg, err := cfg.ConvergenceRadius.LessThanOrEqual(zero)
	if err != nil {
		return Config[Decimal]{}, fmt.Errorf("checking if convergence radius is less than or equal to zero: %w", err)
	}

	if neg {
		return Config[Decimal]{}, ErrConfigConvergenceRadiusPositive
	}

	if cfg.MaxTermsCache == 0 {
		cfg.MaxTermsCache = DefaultMaxTermsCache
	}

	return cfg, nil
}
