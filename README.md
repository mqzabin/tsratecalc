# Taylor Series Rate Calculator

(This repository is a work in progress)

This module provides a calculator to estimate the n-th root of a number using the Taylor Series method with high performance.

The intent is to provide a way to convert a large interval interest rate (e.g. 1 year), to a smaller interval interest rate (e.g. 1 day), by using the n-th root of the large interval interest rate.

**Why it is a "Rate" calculator?**
> The input provided to the `ComputeRate` should be near 0, like common interest rates. When the value is far from 0, the method will converge with more iterations or not converge at all.

# Roadmap
- [x] Implement fuzzy tests using the generic `"shopspring/decimal".Decimal.Pow` method as a source of truth.
- [x] Implement a benchmark to compare the module performance with the generic `"shopspring/decimal".Decimal.Pow` method.
- [x] Implement shopspring/decimal interface adapters.
- [x] Implement convergence radius analysis. 
- [ ] Publish a blog post explaining the Taylor Series method and the module implementation.

# Current benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/mqzabin/tsratecalc/shopspring
cpu: 12th Gen Intel(R) Core(TM) i7-1265U
BenchmarkCalculator_ComputeRate/tsratecalc-12              77008             15009 ns/op           16409 B/op        435 allocs/op
BenchmarkCalculator_ComputeRate/shopspring-12              20716             58773 ns/op           72386 B/op        780 allocs/op
PASS
ok      github.com/mqzabin/tsratecalc/shopspring        3.133s
```