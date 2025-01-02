package tsratecalc

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

type Config struct {
	Root          uint64
	Precision     uint32
	NewFromInt    func(n int64) (Decimal, error)
	MaxIterations uint32
}

func validateConfig(cfg Config) error {
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
