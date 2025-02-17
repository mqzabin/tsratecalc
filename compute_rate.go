package tsratecalc

import "fmt"

var ErrRateOutsideConvergenceBoundaries = fmt.Errorf("rate is outside convergence boundaries")

// ComputeRate receives a rate value and returns "(1+rate)^(1/root) - 1" using a Taylor Series expansion around rate=0.
// The root is defined in the calculator Config.
//
// The rate value should fall within the Config.ConvergenceRadius interval, around rate=0,
// otherwise ErrRateOutsideConvergenceBoundaries will be returned.
//
// It will return ConvergenceError if the desired precision is not achieved after the maximum number of iterations.
func (c *Calculator[Decimal]) ComputeRate(rate Decimal) (Decimal, error) {
	err := c.validateConvergence(rate)
	if err != nil {
		return c.zero, fmt.Errorf("validating boundaries: %w", err)
	}

	var (
		res = c.zero
		// lastError stores the last computed term. It's used to detail the error message, if it happens.
		lastError = c.zero

		variableComponent = c.one
	)

	// Will loop until what happens first:
	// - the desired precision is achieved.
	// - the maximum number of iterations is achieved.
	for n := uint64(1); n < uint64(len(c.taylorTerms)); n++ {
		// variableComponent is rate^n
		variableComponent, err = variableComponent.Mul(rate)
		if err != nil {
			return c.zero, fmt.Errorf("computing rate^%d: %w", n, err)
		}

		currentTermValue, err := c.taylorTerms[n-1].Mul(variableComponent)
		if err != nil {
			return c.zero, fmt.Errorf("computing current taylor term: %w", err)
		}

		// Error checking
		var shouldStop bool
		{
			currentError := currentTermValue

			currentErrorAbs, err := currentError.Abs()
			if err != nil {
				return c.zero, fmt.Errorf("computing taylor aproximation error absolute value: %w", err)
			}

			b, err := currentErrorAbs.LessThanOrEqual(c.maxError)
			if err != nil {
				return c.zero, fmt.Errorf("checking if current error is less than max error: %w", err)
			}

			lastError = currentErrorAbs
			shouldStop = b
		}

		// Adding the new term to the result.

		res, err = res.Add(currentTermValue)
		if err != nil {
			return c.zero, fmt.Errorf("adding current term to result: %w", err)
		}

		if shouldStop {
			res, err = res.Truncate(c.precision)
			if err != nil {
				return c.zero, fmt.Errorf("rounding final result: %w", err)
			}

			return res, nil
		}
	}

	// The loop has ended due to the maximum number of iterations being achieved.
	return c.zero, &ConvergenceError[Decimal]{
		Precision:     c.precision,
		Rate:          rate,
		Iterations:    len(c.taylorTerms),
		LastError:     lastError,
		PartialResult: res,
	}
}

func (c *Calculator[Decimal]) validateConvergence(rate Decimal) error {
	outOfRange, err := rate.LessThanOrEqual(c.convergenceLowerBoundary)
	if err != nil {
		return fmt.Errorf("comparing rate with lower convergence boundary: %w", err)
	}

	if outOfRange {
		return fmt.Errorf("%w: lower boundary is '%s' and rate to compute is '%s'", ErrRateOutsideConvergenceBoundaries, c.convergenceLowerBoundary.String(), rate.String())
	}

	insideRange, err := rate.LessThanOrEqual(c.convergenceUpperBoundary)
	if err != nil {
		return fmt.Errorf("comparing rate with upper convergence boundary: %w", err)
	}

	if !insideRange {
		return fmt.Errorf("%w: upper boundary is '%s' and rate to compute is '%s'", ErrRateOutsideConvergenceBoundaries, c.convergenceUpperBoundary.String(), rate.String())
	}

	return nil
}
