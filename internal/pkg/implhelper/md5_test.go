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
			md5sum:   "8f9270b235499967e357174ffd89de6f",
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
