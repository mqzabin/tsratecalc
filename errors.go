package tsratecalc

import "fmt"

type ConvergenceError struct {
	Precision     uint32
	Rate          Decimal
	Iterations    uint32
	LastError     Decimal
	PartialResult Decimal
}

func (e *ConvergenceError) Error() string {
	return fmt.Sprintf(
		"rate '%s' could not converge to %d digits of precision, it converged to '%s' with %d iterations, last approximation error was '%s'",
		e.Rate.String(),
		e.Precision,
		e.PartialResult.String(),
		e.Iterations,
		e.LastError.String(),
	)
}
