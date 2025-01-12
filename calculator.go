package tsratecalc

import "fmt"

// Calculator is a calculator to convert an annual interest rate to a daily interest rate.
// Use the NewCalculator function to create it.
type Calculator[Decimal Operator[Decimal]] struct {
	precision uint64
	// maxError is the maximum value for the error on calculations. Its value is 2*10^(-(precision+1)).
	maxError Decimal
	// taylorTerms is an in-memory cache for the Taylor series terms constant multipliers.
	taylorTerms []Decimal

	zero Decimal
	one  Decimal

	convergenceUpperBoundary Decimal
	convergenceLowerBoundary Decimal
}

// NewCalculator receives a precision value and a day count convention to return a Calculator.
func NewCalculator[Decimal Operator[Decimal]](cfg Config[Decimal]) (*Calculator[Decimal], error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	root, err := cfg.NewFromInt(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("creating root decimal from integer %d: %w", cfg.Root, err)
	}

	maxError, err := computeMaxError(cfg.Precision, cfg.NewFromInt)
	if err != nil {
		return nil, fmt.Errorf("computing max error: %w", err)
	}

	taylorTerms, err := computeTaylorTermsCache(root, cfg.ConvergenceRadius, maxError, cfg.Precision, cfg.NewFromInt)
	if err != nil {
		return nil, fmt.Errorf("computing taylor terms cache: %w", err)
	}

	zero, err := cfg.NewFromInt(0)
	if err != nil {
		return nil, fmt.Errorf("creating '0' decimal: %w", err)
	}

	one, err := cfg.NewFromInt(1)
	if err != nil {
		return nil, fmt.Errorf("creating '1' decimal: %w", err)
	}

	upperConvergenceBoundary := cfg.ConvergenceRadius

	lowerConvergenceBoundary, err := zero.Sub(cfg.ConvergenceRadius)
	if err != nil {
		return nil, fmt.Errorf("getting lower convergence boundary: %w", err)
	}

	return &Calculator[Decimal]{
		precision:                cfg.Precision,
		maxError:                 maxError,
		taylorTerms:              taylorTerms,
		zero:                     zero,
		one:                      one,
		convergenceUpperBoundary: upperConvergenceBoundary,
		convergenceLowerBoundary: lowerConvergenceBoundary,
	}, nil
}

func (c *Calculator[Decimal]) TermsCacheLen() int {
	return len(c.taylorTerms)
}
