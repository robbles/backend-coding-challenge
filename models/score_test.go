package models

import (
	"testing"
)

func TestRelativeLengthScorer_Score(t *testing.T) {
	tests := map[string]struct {
		query     string
		locations []Location
		location  Location
		expected  float64
	}{
		"matching length": {
			"ABC",
			[]Location{{Name: "ABC"}},
			Location{Name: "ABC"},
			1.0,
		},
		"longer result": {
			"ABC",
			[]Location{{Name: "ABCD"}},
			Location{Name: "ABCD"},
			InverseLengthScore(1),
		},
		"empty string": {
			"",
			[]Location{{Name: ""}},
			Location{Name: ""},
			1.0,
		},
		"multiple results don't affect score": {
			"",
			[]Location{{Name: "ABC"}, {Name: "DEF"}, {Name: "GHI"}},
			Location{Name: "DEF"},
			InverseLengthScore(3),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			scorer := NewRelativeLengthScorer(tt.locations, tt.query)

			if actual := scorer.Score(tt.location); actual != tt.expected {
				t.Errorf("%f", actual)
				t.Errorf("%v != %v", actual, tt.expected)
			}
		})
	}
}
