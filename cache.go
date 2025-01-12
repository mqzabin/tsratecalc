package tsratecalc

import (
	"fmt"
)

// computeTaylorTermsCache creates an in-memory cache of the constant part of the Taylor series terms.
// It will compute all the "terms"-first terms for the provided day count convention.
func computeTaylorTermsCache[Decimal Operator[Decimal]](
	root Decimal,
	convergenceRadius Decimal,
	maxError Decimal,
	precision uint64,
	newFromInt func(n uint64) (Decimal, error),
) ([]Decimal, error) {
	zero, err := newFromInt(0)
	if err != nil {
		return nil, fmt.Errorf("creating '0' decimal: %w", err)
	}

	one, err := newFromInt(1)
	if err != nil {
		return nil, fmt.Errorf("creating '1' decimal: %w", err)
	}

	upperConvergenceBoundary := convergenceRadius

	lowerConvergenceBoundary, err := zero.Sub(convergenceRadius)
	if err != nil {
		return nil, fmt.Errorf("getting lower convergence boundary: %w", err)
	}

	var terms []Decimal

	// Auxiliary accumulators
	var (
		nextDerivativeTermAcc       = zero
		derivativeTermAcc           = one
		factorialTermAcc            = one
		lowerBoundVariableComponent = one
		lastLowerBoundaryError      = zero
		upperBoundVariableComponent = one
		lastUpperBoundaryError      = zero
	)

	// TODO: Remove magic number
	for n := uint64(1); n < 30000; n++ {
		nDecimal, err := newFromInt(n)
		if err != nil {
			return nil, fmt.Errorf("creating '%d' decimal from integer: %w", n, err)
		}

		factorialTermAcc, err = factorialTermAcc.Mul(nDecimal)
		if err != nil {
			return nil, fmt.Errorf("computing factorial term: %w", err)
		}

		derivativeTermAcc, err = derivativeTermAcc.DivRound(root, precision+1)
		if err != nil {
			return nil, fmt.Errorf("computing derivative term: %w", err)
		}

		// derivativeTermAcc * (1 - n * root)
		{
			// n * root
			v, err := nDecimal.Mul(root)
			if err != nil {
				return nil, fmt.Errorf("multiplying n by root: %w", err)
			}

			// 1 - n * root
			v, err = one.Sub(v)
			if err != nil {
				return nil, fmt.Errorf("computing 1 - n*root: %w", err)
			}

			// derivativeTermAcc * (1 - n * root)
			v, err = derivativeTermAcc.Mul(v)
			if err != nil {
				return nil, fmt.Errorf("multiplying derivative term by 1 - n*root: %w", err)
			}

			nextDerivativeTermAcc = v
		}

		// derivativeTermAcc / n!

		term, err := derivativeTermAcc.DivRound(factorialTermAcc, precision+1)
		if err != nil {
			return nil, fmt.Errorf("computing derivative term divided by factorial term: %w", err)
		}

		truncatedTerm, err := term.Truncate(precision + 1)
		if err != nil {
			return nil, fmt.Errorf("truncating taylor term '%s': %w", term.String(), err)
		}

		terms = append(terms, truncatedTerm)
		derivativeTermAcc = nextDerivativeTermAcc

		// Checking the error on lower convergence boundary
		{
			// (lower bound rate)^n
			lowerBoundVariableComponent, err = lowerBoundVariableComponent.Mul(lowerConvergenceBoundary)
			if err != nil {
				return nil, fmt.Errorf("computing lower convergence rate variable component x^%d: %w", n, err)
			}

			// Multiplying term by variable component.
			lowerBoundaryError, err := truncatedTerm.Mul(lowerBoundVariableComponent)
			if err != nil {
				return nil, fmt.Errorf("computing lower boundary error on iteration %d: %w", n, err)
			}

			lowerBoundaryError, err = lowerBoundaryError.Abs()
			if err != nil {
				return nil, fmt.Errorf("computing lower boundary absolute value error on iteration %d: %w", n, err)
			}

			// Should not check first iteration.
			if n > 1 {
				converging, err := lowerBoundaryError.LessThanOrEqual(lastLowerBoundaryError)
				if err != nil {
					return nil, fmt.Errorf("comparing lower boundary error with the last seen on iteration %d: %w", n, err)
				}

				if !converging {
					return nil, fmt.Errorf("lower convergence boundary is diverging (%d taylor term): %w", n, err)
				}

			}

			lastLowerBoundaryError = lowerBoundaryError
		}

		// Checking the error on upper convergence boundary
		{
			// (upper bound rate)^n
			upperBoundVariableComponent, err = upperBoundVariableComponent.Mul(upperConvergenceBoundary)
			if err != nil {
				return nil, fmt.Errorf("computing upper convergence rate variable component x^%d: %w", n, err)
			}

			// Multiplying term by variable component.
			upperBoundaryError, err := truncatedTerm.Mul(upperBoundVariableComponent)
			if err != nil {
				return nil, fmt.Errorf("computing upper boundary error on iteration %d: %w", n, err)
			}

			upperBoundaryError, err = upperBoundaryError.Abs()
			if err != nil {
				return nil, fmt.Errorf("computing upper boundary absolute value error on iteration %d: %w", n, err)
			}

			// Should not check first iteration.
			if n > 1 {
				converging, err := upperBoundaryError.LessThanOrEqual(lastUpperBoundaryError)
				if err != nil {
					return nil, fmt.Errorf("comparing upper boundary error with the last seen on iteration %d: %w", n, err)
				}

				if !converging {
					return nil, fmt.Errorf("upper convergence boundary is diverging (%d taylor term): %w", n, err)
				}
			}

			lastUpperBoundaryError = upperBoundaryError
		}

		if n == 1 {
			continue
		}

		// Checking if the function should stop generating new terms by comparing it to lower boundary error.
		shouldStop, err := lastLowerBoundaryError.LessThanOrEqual(maxError)
		if err != nil {
			return nil, fmt.Errorf("checking if lower boundary error is less than max error: %w", err)
		}

		if !shouldStop {
			continue
		}

		// Checking if the function should stop generating new terms by comparing to upper boundary error.
		shouldStop, err = lastUpperBoundaryError.LessThanOrEqual(maxError)
		if err != nil {
			return nil, fmt.Errorf("checking if upper boundary error is less than max error: %w", err)
		}

		if !shouldStop {
			continue
		}

		break
	}

	return terms, nil
}

func computeMaxError[Decimal Operator[Decimal]](
	precision uint64,
	newFromInt func(n uint64) (Decimal, error),
) (Decimal, error) {
	one, err := newFromInt(1)
	if err != nil {
		var zero Decimal

		return zero, fmt.Errorf("creating '1' decimal: %w", err)
	}

	two, err := newFromInt(2)
	if err != nil {
		var zero Decimal

		return zero, fmt.Errorf("creating '2' decimal: %w", err)
	}

	ten, err := newFromInt(10)
	if err != nil {
		var zero Decimal

		return zero, fmt.Errorf("creating '10' decimal: %w", err)
	}

	// maxError is the maximum value for the error on calculations. Its value is 2*10^(-(precision+1)).

	// 10^(precision+1)
	maxError, err := ten.PowInt(precision)
	if err != nil {
		var zeroValue Decimal

		return zeroValue, fmt.Errorf("raising 10 to the power of precision+1: %w", err)
	}

	// 2*10^(precision+1)
	maxError, err = two.Mul(maxError)
	if err != nil {
		var zeroValue Decimal

		return zeroValue, fmt.Errorf("multiplying max error by 2: %w", err)
	}

	// 1/(2*10^(precision+1))
	maxError, err = one.DivRound(maxError, precision+1)
	if err != nil {
		var zeroValue Decimal

		return zeroValue, fmt.Errorf("dividing 1 by max error: %w", err)
	}

	return maxError, nil
}
