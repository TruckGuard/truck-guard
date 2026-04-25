package utils

// min3 returns the minimum of three integers.
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// LevenshteinDistance calculates the minimum number of single-character edits
// (insertions, deletions, or substitutions) required to change string s1 into s2.
// It supports UTF-8 runes (e.g. Cyrillic vs Latin characters).
func LevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	r1 := []rune(s1)
	r2 := []rune(s2)
	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create a matrix of size (len1+1) x (len2+1)
	// We use two rows to optimize memory allocation (O(min(len1, len2)))
	if len1 > len2 {
		r1, r2 = r2, r1
		len1, len2 = len2, len1
	}

	row1 := make([]int, len1+1)
	row2 := make([]int, len1+1)

	for i := 0; i <= len1; i++ {
		row1[i] = i
	}

	for i := 0; i < len2; i++ {
		row2[0] = i + 1
		for j := 0; j < len1; j++ {
			cost := 1
			if r1[j] == r2[i] {
				cost = 0
			}
			row2[j+1] = min3(
				row1[j+1]+1, // deletion
				row2[j]+1,   // insertion
				row1[j]+cost, // substitution
			)
		}
		copy(row1, row2)
	}

	return row1[len1]
}
