package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplitStringByChars(t *testing.T) {
	var testCases = []struct {
		given string
		size  int
		exp   []string
	}{
		{"", 5, []string{"", ""}},
		{"こんにちは世界", 5, []string{"こんにちは", "世界"}},
		{"こんにちは世界", 10, []string{"こんにちは世界", ""}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("s=%s, size=%d", tc.given, tc.size), func(t *testing.T) {
			t.Parallel()
			s1, s2 := SplitStringByChars(tc.given, tc.size)
			if !reflect.DeepEqual([]string{s1, s2}, tc.exp) {
				t.Errorf("expected: %s, got: %s", tc.exp, []string{s1, s2})
			}
		})
	}
}
