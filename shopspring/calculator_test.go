package shopspring_test

import (
	"testing"

	"github.com/mqzabin/fuzzdecimal"
	"github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc/shopspring"
)

func FuzzComputeRateShopspring(f *testing.F) {
	const (
		maxIterations     = 100
		resultPrecision   = 30
		divisionPrecision = 30
		root              = 252
	)

	cfg := shopspring.Config{
		Root:            root,
		MaxIterations:   maxIterations,
		ResultPrecision: resultPrecision,
	}

	calc, err := shopspring.NewCalculator(cfg)
	if err != nil {
		f.Fatalf("NewCalculator: %v", err)
	}

	parseDecimal := func(t *fuzzdecimal.T, s string) (decimal.Decimal, error) {
		t.Helper()

		return decimal.NewFromString(s)
	}

	one := decimal.NewFromInt(1)
	refExponent := one.DivRound(decimal.NewFromInt(root), divisionPrecision)

	fuzzdecimal.Fuzz(f, 1, func(t *fuzzdecimal.T) {
		fuzzdecimal.AsDecimalComparison1(t, "ComputeRate", parseDecimal, parseDecimal,
			func(t *fuzzdecimal.T, x1 decimal.Decimal) (string, error) {
				t.Helper()

				return x1.Add(one).Pow(refExponent).Sub(one).Truncate(resultPrecision).StringFixed(resultPrecision), nil
			},
			func(t *fuzzdecimal.T, x1 decimal.Decimal) string {
				res, err := calc.ComputeRate(x1)
				if err != nil {
					t.Fatalf("ComputeRate: %v", err)
				}

				t.Log(x1.String())

				return res.StringFixed(resultPrecision)
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithMaxSignificantDigits(resultPrecision),
		fuzzdecimal.WithMaxDecimalPlaces(resultPrecision),
		fuzzdecimal.WithUnsigned(),
	))
}
