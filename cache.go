package tsratecalc

import (
	"fmt"
)

func computeIntToDecimalCache[Decimal Operator[Decimal]](newFromInt func(uint64) (Decimal, error), maxValue uint64) ([]Decimal, error) {
	cacheSize := maxValue + 1

	intToDecimal := make([]Decimal, cacheSize)
	for i := range uint64(len(intToDecimal)) {
		dec, err := newFromInt(i)
		if err != nil {
			return nil, fmt.Errorf("creating '%d' decimal from integer: %w", i, err)
		}

		intToDecimal[i] = dec
	}

	return intToDecimal, nil
}

func computeMaxError[Decimal Operator[Decimal]](intToDecimal []Decimal, precision uint64) (Decimal, error) {
	one := intToDecimal[1]
	two := intToDecimal[2]
	ten := intToDecimal[10]

	// maxError is the maximum value for the error on calculations. Its value is 2*10^(-(precision+1)).

	// 10^(precision+1)
	maxError, err := ten.PowInt(precision)
	if err != nil {
		var zeroValue Decimal

		return zeroValue, fmt.Errorf("raising 10 to the power of precision+1: %w", err)
	}

	// 2*10^(precision+1)
	{
		v, err := two.Mul(maxError)
		if err != nil {
			var zeroValue Decimal

			return zeroValue, fmt.Errorf("multiplying max error by 2: %w", err)
		}

		maxError = v
	}

	// 1/(2*10^(precision+1))
	{
		v, err := one.DivRound(maxError, precision+1)
		if err != nil {
			var zeroValue Decimal

			return zeroValue, fmt.Errorf("dividing 1 by max error: %w", err)
		}

		maxError = v
	}

	return maxError, nil
}

// computeTaylorTermsCache creates an in-memory cache of the constant part of the Taylor series terms.
// It will compute all the "terms"-first terms for the provided day count convention.
func computeTaylorTermsCache[Decimal Operator[Decimal]](root Decimal, intToDecimal []Decimal, precision, terms uint64) ([]Decimal, error) {
	cache := make([]Decimal, terms)

	one := intToDecimal[1]

	// First term around 1 is 1.
	cache[0] = intToDecimal[1]

	// Auxiliary accumulators
	var nextDerivativeTermAcc Decimal
	derivativeTermAcc := one
	factorialTermAcc := one

	for n := uint64(1); n < terms; n++ {
		// n! = n*(n-1)!
		{
			v, err := factorialTermAcc.Mul(intToDecimal[n])
			if err != nil {
				return nil, fmt.Errorf("computing factorial term: %w", err)
			}

			factorialTermAcc = v
		}

		// derivativeTermAcc
		{
			v, err := derivativeTermAcc.DivRound(root, precision+1)
			if err != nil {
				return nil, fmt.Errorf("computing derivative term: %w", err)
			}

			derivativeTermAcc = v
		}

		// derivativeTermAcc * (1 - n * root)
		{
			// n * root
			v, err := intToDecimal[n].Mul(root)
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
		{
			v, err := derivativeTermAcc.DivRound(factorialTermAcc, precision+1)
			if err != nil {
				return nil, fmt.Errorf("computing derivative term divided by factorial term: %w", err)
			}

			cache[n] = v
		}

		derivativeTermAcc = nextDerivativeTermAcc
	}

	return cache, nil
}