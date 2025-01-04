package shopspring_test

import (
	"errors"
	"testing"

	"github.com/mqzabin/tsratecalc/shopspring"
)

func TestCalculator_NewCalculator_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		cfg  shopspring.Config
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := shopspring.NewCalculator(tc.cfg)
			if c == nil {
				t.Errorf("expected non-nil calculator")
			}

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestCalculator_NewCalculator_Failure(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		cfg      shopspring.Config
		expected error
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := shopspring.NewCalculator(tc.cfg)
			if c != nil {
				t.Errorf("expected nil calculator")
			}

			if err == nil {
				t.Errorf("expected error, got nil")
			}

			if !errors.Is(err, tc.expected) {
				t.Errorf("expected error %v, got %v", tc.expected, err)
			}
		})
	}
}
