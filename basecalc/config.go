package basecalc

import (
	"errors"
	"fmt"
)

const (
	minPrecision     = 0
	minRoot          = 2
	minMaxIterations = 1
)

var (
	ErrConfigRootMinValue      = fmt.Errorf("root should be greater than %d", minRoot)
	ErrConfigPrecisionMinValue = fmt.Errorf("precision should be greater than %d", minPrecision)
	ErrConfigMaxIterations     = fmt.Errorf("max iterations should be greater than %d", minMaxIterations)
	ErrConfigNewFromIntIsNil   = errors.New("'decimal from integer' factory should not be nil")
)

type Config[Decimal Operator[Decimal]] struct {
	Root          uint64
	Precision     uint64
	NewFromInt    func(n uint64) (Decimal, error)
	MaxIterations uint64
}

func validateConfig[Decimal Operator[Decimal]](cfg Config[Decimal]) error {
	if cfg.Root < minRoot {
		return ErrConfigRootMinValue
	}

	if cfg.Precision < minPrecision {
		return ErrConfigPrecisionMinValue
	}

	if cfg.MaxIterations < minMaxIterations {
		return ErrConfigMaxIterations
	}

	if cfg.NewFromInt == nil {
		return ErrConfigNewFromIntIsNil
	}

	return nil
}
