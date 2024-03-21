package internal

import "unicode/utf8"

// SplitStringByChars slices the input string into two parts:
// 1. has exactly size characters
// 2. contains the remaining characters
func SplitStringByChars(s string, size int) (string, string) {
	// If the string is empty or size is 0, return empty strings
	if s == "" || size == 0 {
		return "", s
	}
	// Initialize variables to track the indices for slicing
	var i, runeCount int
	// Iterate over the string
	for i < len(s) && runeCount < size {
		// Decode the next rune
		// NOTE: i is incremented by the byte size of runes so slicing
		// the string s is safe without breaking multibyte chars
		_, runeSize := utf8.DecodeRuneInString(s[i:])
		// Increase rune count and move to the next rune
		runeCount++
		i += runeSize
	}
	// Return the first part of the string and the remaining part
	return s[:i], s[i:]
}
