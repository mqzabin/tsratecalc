package shopspring_test

import (
	"fmt"
	"testing"

	"github.com/mqzabin/fuzzdecimal"
	"github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc/shopspring"
)

func TestNewCalculator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config shopspring.Config
	}{
		{
			name: "",
			config: shopspring.Config{
				Root:              252,
				Precision:         30,
				ConvergenceRadius: decimal.New(9, -1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			calc, err := shopspring.NewCalculator(tc.config)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			fmt.Println(calc.TermsCacheLen())
		})
	}
}

func FuzzComputeRateShopspring(f *testing.F) {
	const (
		resultPrecision   = 30
		divisionPrecision = 30
		root              = 252
	)

	cfg := shopspring.Config{
		Root:              root,
		Precision:         resultPrecision,
		ConvergenceRadius: decimal.New(9, -1), // 0.9
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

				return res.StringFixed(resultPrecision)
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithMaxSignificantDigits(resultPrecision),
		fuzzdecimal.WithMaxDecimalPlaces(resultPrecision),
		fuzzdecimal.WithUnsigned(),
	))
}
