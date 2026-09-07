package contracts

import (
	"math"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

func TestMoneyUsesExactMinorUnits(t *testing.T) {
	for _, tc := range []struct {
		value string
		cents int64
	}{
		{"0.01", 1}, {"1.01", 101}, {"1.10", 110}, {"99.9", 9990},
		{"0001.20", 120}, {" 2500.50 ", 250050},
		{"92233720368547758.07", math.MaxInt64},
	} {
		t.Run(tc.value, func(t *testing.T) {
			got, err := parseContractValueCents(tc.value)
			if err != nil || got != tc.cents {
				t.Fatalf("got %d, %v; want %d", got, err, tc.cents)
			}
			form := contractFormFromModel(models.Contract{ValueCents: &got})
			back, err := parseContractValueCents(form.Value)
			if err != nil || back != got {
				t.Fatalf("round trip %q: %d, %v", form.Value, back, err)
			}
		})
	}
}

func TestMoneyRejectsNonFiniteFractionalAndOverflowValues(t *testing.T) {
	for _, value := range []string{"NaN", "Inf", "+Inf", "-Inf", "1e3", "-0", "+1", "-1", "1.001", "1.", ".50", "1..2", "1,000", "１２", "92233720368547758.08", "999999999999999999999"} {
		t.Run(value, func(t *testing.T) {
			if _, err := parseContractValueCents(value); err == nil {
				t.Fatalf("accepted %q", value)
			}
		})
	}
}

func FuzzContractValueCents(f *testing.F) {
	for _, seed := range []string{"0", "1.01", "NaN", "92233720368547758.07", "1.001", ""} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, value string) {
		cents, err := parseContractValueCents(value)
		if err != nil {
			return
		}
		if cents < 0 {
			t.Fatalf("negative cents for %q", value)
		}
		form := contractFormFromModel(models.Contract{ValueCents: &cents})
		got, err := parseContractValueCents(form.Value)
		if err != nil || got != cents {
			t.Fatalf("failed round trip for %q", value)
		}
	})
}
