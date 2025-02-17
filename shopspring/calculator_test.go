package shopspring_test

import (
	"testing"

	"github.com/mqzabin/fuzzdecimal"
	"github.com/shopspring/decimal"

	"github.com/mqzabin/tsratecalc/shopspring"
)

func BenchmarkCalculator_ComputeRate_30Digits(b *testing.B) {
	const (
		resultPrecision   = 30
		divisionPrecision = 30
		root              = 252
	)

	rate := decimal.New(1, -1) // 10%

	b.ReportAllocs()

	b.Run("tsratecalc", func(b *testing.B) {
		cfg := shopspring.Config{
			Root:              root,
			Precision:         resultPrecision,
			ConvergenceRadius: decimal.New(9, -1),
		}

		calc, err := shopspring.NewCalculator(cfg)
		if err != nil {
			b.Fatalf("NewCalculator: %v", err)
		}

		var avoidOptimizations decimal.Decimal

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			avoidOptimizations, _ = calc.ComputeRate(rate)
		}

		if avoidOptimizations.IsZero() {
			b.Fatalf("unexpected zero result")
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		one := decimal.NewFromInt(1)
		refExponent := one.DivRound(decimal.NewFromInt(root), divisionPrecision)

		var avoidOptimizations decimal.Decimal

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			avoidOptimizations = rate.Add(one).Pow(refExponent).Sub(one).Truncate(resultPrecision)
		}

		if avoidOptimizations.IsZero() {
			b.Fatalf("unexpected zero result")
		}
	})
}

func BenchmarkCalculator_ComputeRate_10Digits(b *testing.B) {
	const (
		resultPrecision   = 10
		divisionPrecision = 30
		root              = 252
	)

	rate := decimal.New(1, -1) // 10%

	b.ReportAllocs()

	b.Run("tsratecalc", func(b *testing.B) {
		cfg := shopspring.Config{
			Root:              root,
			Precision:         resultPrecision,
			ConvergenceRadius: decimal.New(9, -1),
		}

		calc, err := shopspring.NewCalculator(cfg)
		if err != nil {
			b.Fatalf("NewCalculator: %v", err)
		}

		var avoidOptimizations decimal.Decimal

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			avoidOptimizations, _ = calc.ComputeRate(rate)
		}

		if avoidOptimizations.IsZero() {
			b.Fatalf("unexpected zero result")
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		one := decimal.NewFromInt(1)
		refExponent := one.DivRound(decimal.NewFromInt(root), divisionPrecision)

		var avoidOptimizations decimal.Decimal

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			avoidOptimizations = rate.Add(one).Pow(refExponent).Sub(one).Truncate(resultPrecision)
		}

		if avoidOptimizations.IsZero() {
			b.Fatalf("unexpected zero result")
		}
	})
}

func TestNewCalculator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		config            shopspring.Config
		wantTermsCacheLen int
	}{
		{
			name: "30 digits with 0.9 convergence radius",
			config: shopspring.Config{
				Root:              252,
				Precision:         30,
				ConvergenceRadius: decimal.New(9, -1),
			},
			wantTermsCacheLen: 550,
		},
		{
			name: "30 digits with 0.8 convergence radius",
			config: shopspring.Config{
				Root:              252,
				Precision:         30,
				ConvergenceRadius: decimal.New(8, -1),
			},
			wantTermsCacheLen: 263,
		},
		{
			name: "10 digits with 0.9 convergence radius",
			config: shopspring.Config{
				Root:              252,
				Precision:         10,
				ConvergenceRadius: decimal.New(9, -1),
			},
			wantTermsCacheLen: 127,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			calc, err := shopspring.NewCalculator(tc.config)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if calc == nil {
				t.Fatalf("unexpected nil calculator")
			}

			if tc.wantTermsCacheLen != calc.TermsCacheLen() {
				t.Fatalf("unexpected terms cache length: got %d, want %d", calc.TermsCacheLen(), tc.wantTermsCacheLen)
			}
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
