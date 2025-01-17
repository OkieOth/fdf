package implhelper_test

import (
	"testing"

	"github.com/okieoth/fdf/internal/pkg/implhelper"
)

func TestGetMd5(t *testing.T) {
	type TestCases struct {
		filePath string
		md5sum   string
	}
	testCases := []TestCases{
		{
			filePath: "../../../.gitignore",
			md5sum:   "a945dada977c23b702324087618ea4ee",
		},
		{
			filePath: "../../../LICENSE",
			md5sum:   "39dce100a9679cf1aee38c6ea4357fc6",
		},
	}

	for _, c := range testCases {
		s, e := implhelper.GetMd5(c.filePath)
		if e != nil {
			t.Errorf("Error in getting md5 - file: %s, err: %s", c.filePath, e)
		} else {
			if s != c.md5sum {
				t.Errorf("Wrong md5 sum - expected: %s, got: %s", c.md5sum, s)
			}
		}
	}
}

func TestGetCommonPrefix(t *testing.T) {
	// Happy cases
	tests := []struct {
		s1       string
		s2       string
		expected string
	}{
		{"hello", "helicopter", "hel"},
		{"abcdef", "abcxyz", "abc"},
		{"prefix", "prefixmatch", "prefix"},
		{"same", "same", "same"},
		{"short", "shorter", "short"},
	}

	for _, tt := range tests {
		t.Run("Happy case: "+tt.s1+" & "+tt.s2, func(t *testing.T) {
			result := implhelper.GetCommonPrefix(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}

	// Unhappy cases
	unhappyTests := []struct {
		s1       string
		s2       string
		expected string
	}{
		{"abc", "xyz", ""},
		{"different", "prefix", ""},
		{"", "nonempty", ""},
		{"nonempty", "", ""},
		{"", "", ""},
	}

	for _, tt := range unhappyTests {
		t.Run("Unhappy case: "+tt.s1+" & "+tt.s2, func(t *testing.T) {
			result := implhelper.GetCommonPrefix(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAdjustCommonPrefix(t *testing.T) {
	// Happy Cases
	t.Run("CommonPrefix", func(t *testing.T) {
		prefix := "hello"
		s := "helloworld"
		expected := "hello"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("ExactMatch", func(t *testing.T) {
		prefix := "test"
		s := "test"
		expected := "test"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("ShorterPrefix", func(t *testing.T) {
		prefix := "go"
		s := "golang"
		expected := "go"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	// Unhappy Cases
	t.Run("NoCommonPrefix", func(t *testing.T) {
		prefix := "abc"
		s := "xyz"
		expected := ""
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("EmptyPrefix", func(t *testing.T) {
		prefix := ""
		s := "anything"
		expected := ""
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("EmptyString", func(t *testing.T) {
		prefix := "something"
		s := ""
		expected := ""
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("UnicodeCharacters", func(t *testing.T) {
		prefix := "你好世界"
		s := "你好golang"
		expected := "你好"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		prefix := "!@#$"
		s := "!@#abc"
		expected := "!@#"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("PrefixLongerThanString", func(t *testing.T) {
		prefix := "longprefix"
		s := "long"
		expected := "long"
		result := implhelper.AdjustCommonPrefix(prefix, s)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}
