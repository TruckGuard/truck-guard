package utils

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"AA1234BB", "AA1234BB", 0},
		{"AA1234BB", "AA1234B8", 1}, // B -> 8
		{"AA1234BB", "AA12340B", 1}, // B -> 0
		{"AA1234BB", "AA1234", 2},   // deletion of BB
		{"AA1234", "AA1234BB", 2},   // insertion of BB
		{"AA1234BB", "AА1234BB", 1}, // Latin A vs Cyrillic А (different runes)
	}

	for _, tt := range tests {
		t.Run(tt.s1+"_"+tt.s2, func(t *testing.T) {
			actual := LevenshteinDistance(tt.s1, tt.s2)
			if actual != tt.expected {
				t.Errorf("LevenshteinDistance(%q, %q) = %d; want %d", tt.s1, tt.s2, actual, tt.expected)
			}
		})
	}
}
