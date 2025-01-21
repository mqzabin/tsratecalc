package tsratecalc

import "fmt"

// ConvergenceError is an error type for when the rate value could not converge to the desired precision.
type ConvergenceError[Decimal Operator[Decimal]] struct {
	// Precision is the desired number of decimal places to consider in the calculations.
	Precision uint64
	// Rate is the rate value that could not converge.
	Rate Decimal
	// Iterations is the number of iterations that were performed.
	Iterations int
	// LastError is the last approximation error.
	LastError Decimal
	// PartialResult is the partial result of the calculation.
	PartialResult Decimal
}

func (e *ConvergenceError[Decimal]) Error() string {
	return fmt.Sprintf(
		"rate '%s' could not converge to %d digits of precision, it converged to '%s' with %d iterations, last approximation error was '%s'",
		e.Rate.String(),
		e.Precision,
		e.PartialResult.String(),
		e.Iterations,
		e.LastError.String(),
	)
}
