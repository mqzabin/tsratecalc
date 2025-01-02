package tsratecalc

import "fmt"

// Calculator is a calculator to convert an annual interest rate to a daily interest rate.
// Use the NewCalculator function to create it.
type Calculator struct {
	precision uint32
	// maxError is the maximum value for the error on calculations. Its value is 2*10^(-(precision+1)).
	maxError Decimal
	// taylorTerms is an in-memory cache for the Taylor series terms constant multipliers.
	taylorTerms []Decimal
	// intToDecimal is an in-memory cache for the integers used in Taylor series calculation.
	intToDecimal []Decimal
	// maxIterations represents the Taylor series' maximum number of iterations (or terms).
	maxIterations uint32

	newFromInt func(n int64) (Decimal, error)
}

// NewCalculator receives a precision value and a day count convention to return a Calculator.
func NewCalculator(cfg Config) (*Calculator, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	intToDecimal, err := computeIntToDecimalCache(cfg.NewFromInt, cfg.MaxIterations)
	if err != nil {
		return nil, fmt.Errorf("building int to decimal cache: %w", err)
	}

	root, err := cfg.NewFromInt(int64(cfg.Root))
	if err != nil {
		return nil, fmt.Errorf("creating root decimal from integer %d: %w", cfg.Root, err)
	}

	maxError, err := computeMaxError(intToDecimal, cfg.Precision)
	if err != nil {
		return nil, fmt.Errorf("computing max error: %w", err)
	}

	taylorTerms, err := computeTaylorTermsCache(root, intToDecimal, cfg.Precision)
	if err != nil {
		return nil, fmt.Errorf("computing taylor terms cache: %w", err)
	}

	return &Calculator{
		precision:     cfg.Precision,
		maxError:      maxError,
		taylorTerms:   taylorTerms,
		intToDecimal:  intToDecimal,
		maxIterations: cfg.MaxIterations,
		newFromInt:    cfg.NewFromInt,
	}, nil
}

// ComputeRate receives an annual interest rate and returns the daily interest rate.
// The provided rate should belong to the [ -0.8, 0.8 ] interval.
//
// !!ATTENTION!! The rate value should be provided in decimal format, e.g. inform 0.13 for a 13% annual interest rate.
//
// It could return an ErrMaxIterationsAchieved error if the max number of iterations is achieved without convergence.
// If you receive ErrMaxIterationsAchieved it's possible that you've set a too high precision or the rate value
// is outside the [ -0.8, 0.8 ] interval (or near to the boundaries).
func (rc *Calculator) ComputeRate(rate Decimal) (Decimal, error) {
	one := rc.intToDecimal[1]
	zero := rc.intToDecimal[0]

	// Adding 1 to the provided rate, e.g. 0.13 turns into 1.13.
	var shiftedRate Decimal
	{
		v, err := rate.Add(one)
		if err != nil {
			return nil, fmt.Errorf("adding 1 to rate: %w", err)
		}

		shiftedRate = v
	}

	// First Taylor series term.
	res := one

	// lastError stores the last computed term. It's used to detail the error message, if it happens.
	var lastError Decimal

	// Will loop until what happens first:
	// - the desired precision is achieved.
	// - the maximum number of iterations is achieved.
	for iteration := uint32(1); iteration < rc.maxIterations; iteration++ {

		// variableComponent is (shiftedRate - 1)^iteration
		var variableComponent Decimal
		{
			sub, err := shiftedRate.Sub(one)
			if err != nil {
				return nil, fmt.Errorf("computing rate - 1: %w", err)
			}

			// (rate - 1)^iteration
			power, err := sub.PowInt(iteration)
			if err != nil {
				return nil, fmt.Errorf("computing (rate - 1)^n: %w", err)
			}

			variableComponent = power
		}

		var currentTermValue Decimal
		{
			// Gets the cached constant multiplier.
			taylorTerm := rc.taylorTerms[iteration]

			v, err := taylorTerm.Mul(variableComponent)
			if err != nil {
				return nil, fmt.Errorf("computing current taylor term: %w", err)
			}

			currentTermValue = v
		}

		// Error checking
		var shouldStop bool
		{
			currentError := currentTermValue

			currentErrorAbs, err := currentError.Abs()
			if err != nil {
				return nil, fmt.Errorf("computing taylor aproximation error absolute value: %w", err)
			}

			b, err := currentErrorAbs.LessThan(rc.maxError)
			if err != nil {
				return nil, fmt.Errorf("checking if current error is less than max error: %w", err)
			}

			lastError = currentErrorAbs

			shouldStop = b
		}

		if shouldStop {
			v, err := res.Sub(one)
			if err != nil {
				return nil, fmt.Errorf("subtracting 1 from end result: %w", err)
			}

			res = v

			return res, nil
		}

		// Adding the new term to the result.
		{
			v, err := res.Add(currentTermValue)
			if err != nil {
				return nil, fmt.Errorf("adding current term to result: %w", err)
			}

			res = v
		}
	}

	// The loop has ended due to the maximum number of iterations being achieved.
	return zero, &ConvergenceError{
		Precision:     rc.precision,
		Rate:          rate,
		Iterations:    rc.maxIterations,
		LastError:     lastError,
		PartialResult: res,
	}
}
