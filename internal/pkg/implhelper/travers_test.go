package implhelper

import (
	"testing"
)

func TestBlackWhiteListMatch(t *testing.T) {
	tests := []struct {
		name           string
		patternToMatch string
		path           string
		expected       bool
	}{
		// Happy cases
		{
			name:           "Exact match without wildcard",
			patternToMatch: "example.jpeg",
			path:           "example.jpeg",
			expected:       true,
		},
		{
			name:           "Wildcard match at the beginning",
			patternToMatch: "*.jpeg",
			path:           "image.jpeg",
			expected:       true,
		},
		{
			name:           "Wildcard match in the middle",
			patternToMatch: "IMG*.jpeg",
			path:           "IMG123.jpeg",
			expected:       true,
		},
		{
			name:           "Wildcard match with multiple wildcards",
			patternToMatch: "IMG*.j*",
			path:           "IMG123.jpg",
			expected:       true,
		},
		{
			name:           "Substring match without wildcard",
			patternToMatch: "example",
			path:           "example.jpeg",
			expected:       true,
		},

		// Unhappy cases
		{
			name:           "No match without wildcard",
			patternToMatch: "example.jpeg",
			path:           "image.jpeg",
			expected:       false,
		},
		{
			name:           "Wildcard match fails",
			patternToMatch: "IMG*.jpeg",
			path:           "image.jpeg",
			expected:       false,
		},
		{
			name:           "Substring match fails",
			patternToMatch: "example",
			path:           "image.jpeg",
			expected:       false,
		},
		{
			name:           "Wildcard match fails with additional characters",
			patternToMatch: "IMG*.jpeg",
			path:           "IMG123.png",
			expected:       false,
		},
		{
			name:           "Empty pattern",
			patternToMatch: "",
			path:           "image.jpeg",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blackWhileListMatch(tt.patternToMatch, tt.path)
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}

}
