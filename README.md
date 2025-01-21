# Taylor Series Rate Calculator

This module provides a generic calculator to estimate the function:

$$ f(x) = \sqrt[c]{1+x} - 1 \text{, for } c \in \mathbb{N} $$

This function is generally used to convert a large interval interest rate (e.g. 1 year), to a smaller interval interest rate (e.g. 1 day), by using the n-th root of the large interval interest rate.

The Taylor Series is expanded around the point x=0, and you could specify the desired convergence radius.
The calculator will automatically expand the Taylor terms up to the necessary "n" to converge inside the provided radius.

The Taylor Series expansion is given by:

$$ f(x) = \sum_{n=1}^{\infty} \frac{1}{c^{n}n!} \Bigg(\prod_{i=1}^{n-1} (1 - ic) \Bigg)  x^n $$

# Usage

The `tsratecalc` package was not meant to be used directly. It needs an adapter to be used with a specific arbitrary/fixed precision decimal structure.

The required operations are described by the `Operator` interface:

```go
type Operator[Decimal any] interface {
	Mul(n Decimal) (Decimal, error)
	DivRound(n Decimal, places uint64) (Decimal, error)
	Sub(n Decimal) (Decimal, error)
	Add(n Decimal) (Decimal, error)
	Abs() (Decimal, error)
	LessThanOrEqual(n Decimal) (bool, error)
	PowInt(n uint64) (Decimal, error)
	Truncate(places uint64) (Decimal, error)
	String() string
}

```

There are some subpackages that implement the adapters for some decimal types:

- `tsratecalc/shopspring`: Support for the `github.com/shopspring/decimal` package.

# Current benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/mqzabin/tsratecalc/shopspring
cpu: 12th Gen Intel(R) Core(TM) i7-1265U
BenchmarkCalculator_ComputeRate_30Digits/tsratecalc-12             79836             15035 ns/op           16410 B/op        435 allocs/op
BenchmarkCalculator_ComputeRate_30Digits/shopspring-12             20139             57067 ns/op           72386 B/op        780 allocs/op
BenchmarkCalculator_ComputeRate_10Digits/tsratecalc-12            280160              4231 ns/op            4232 B/op        139 allocs/op
BenchmarkCalculator_ComputeRate_10Digits/shopspring-12             19969             57254 ns/op           72445 B/op        780 allocs/op
```