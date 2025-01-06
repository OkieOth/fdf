package cmdimpl_test

import (
	"testing"

	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
)

func TestListImpl(t *testing.T) {
	whiteList := make([]string, 0)
	blackList := make([]string, 0)
	fileRepo := helper.NewFileRepo()
	if err := fileRepo.InitFromSource("/home/eiko/prog/git/fdf/internal", blackList, whiteList); err != nil {
		t.Error(err)
	} else {
		rs := fileRepo.Size()
		if rs == 0 {
			t.Errorf("Seems the fileRepo isn't initialized: %d", rs)
		}
	}
}
