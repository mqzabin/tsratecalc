package tsratecalc

import "fmt"

// Calculator is a calculator for "(1+x)^(1/n)-1", with positive integer n.
// It uses a Taylor series expansion around x=0 to compute the rate value.
//
// It could be used for any arbitrary/fixed precision decimal that implements the Operator interface.
type Calculator[Decimal Operator[Decimal]] struct {
	// precision is the number of decimal places to consider in the calculations.
	precision uint64
	// maxError is the maximum value for the error on calculations. Its value is 2*10^(-(precision+1)).
	maxError Decimal
	// taylorTerms is an in-memory cache for the Taylor series terms constant multipliers.
	taylorTerms []Decimal
	// zero store the zero value for the Decimal type.
	zero Decimal
	// one store the one value for the Decimal type.
	one Decimal
	// convergenceUpperBoundary is the upper boundary for the rate value to be considered inside the convergence radius.
	convergenceUpperBoundary Decimal
	// convergenceLowerBoundary is the lower boundary for the rate value to be considered inside the convergence radius.
	convergenceLowerBoundary Decimal
}

// NewCalculator returns a new Calculator given a Config for a specific Decimal type.
// The Decimal type should implement the Operator interface.
func NewCalculator[Decimal Operator[Decimal]](cfg Config[Decimal]) (*Calculator[Decimal], error) {
	cfg, err := validateConfig(cfg)
	if err != nil {
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

	taylorTerms, err := computeTaylorTermsCache(root, cfg.ConvergenceRadius, cfg.MaxTermsCache, maxError, cfg.Precision, cfg.NewFromInt)
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

// TermsCacheLen returns the number of Taylor terms stored in the calculator's cache.
func (c *Calculator[Decimal]) TermsCacheLen() int {
	return len(c.taylorTerms)
}
