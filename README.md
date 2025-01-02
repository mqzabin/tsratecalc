# Taylor Series Rate Calculator

(This repository is a work in progress)

This module provides a calculator to estimate the n-th root of a number using the Taylor Series method with high performance.

The intent is to provide a way to convert a large interval interest rate (e.g. 1 year), to a smaller interval interest rate (e.g. 1 day), by using the n-th root of the large interval interest rate.

**Why it is a "Rate" calculator?**
> The input provided to the `ComputeRate` should be near 0, like common interest rates. When the value is far from 0, the method will converge with more iterations or not converge at all.

# Roadmap
- [ ] Implement tests (standard and fuzzy) using the generic `"shopspring/decimal".Decimal.Pow` method as a source of truth.
- [ ] Implement a benchmark to compare the module performance with the generic `"shopspring/decimal".Decimal.Pow` method.
- [ ] Implement shopspring/decimal interface adapters.
- [ ] Publish a blog post explaining the Taylor Series method and the module implementation.
- [ ] Implement convergence radius analysis. 