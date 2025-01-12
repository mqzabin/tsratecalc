package tsratecalc

import (
	"errors"
	"fmt"
)

const (
	minPrecision = 0
	minRoot      = 2
)

var (
	ErrConfigRootMinValue              = fmt.Errorf("root should be greater than %d", minRoot)
	ErrConfigPrecisionMinValue         = fmt.Errorf("precision should be greater than %d", minPrecision)
	ErrConfigNewFromIntIsNil           = errors.New("'decimal from integer' factory should not be nil")
	ErrConfigConvergenceRadiusPositive = errors.New("convergence radius must be positive")
)

type Config[Decimal Operator[Decimal]] struct {
	Root              uint64
	Precision         uint64
	NewFromInt        func(n uint64) (Decimal, error)
	ConvergenceRadius Decimal
}

func validateConfig[Decimal Operator[Decimal]](cfg Config[Decimal]) error {
	if cfg.Root < minRoot {
		return ErrConfigRootMinValue
	}

	if cfg.Precision < minPrecision {
		return ErrConfigPrecisionMinValue
	}

	if cfg.NewFromInt == nil {
		return ErrConfigNewFromIntIsNil
	}

	zero, err := cfg.NewFromInt(0)
	if err != nil {
		return fmt.Errorf("creating '0' decimal: %w", err)
	}

	neg, err := cfg.ConvergenceRadius.LessThanOrEqual(zero)
	if err != nil {
		return fmt.Errorf("checking if convergence radius is less than or equal to zero: %w", err)
	}

	if neg {
		return ErrConfigConvergenceRadiusPositive
	}

	return nil
}
