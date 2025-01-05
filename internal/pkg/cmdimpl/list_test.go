package cmdimpl_test

import (
	"fmt"
	"testing"

	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
)

func TestListImpl(t *testing.T) {
	dummy := make([]string, 0)
	fileRepo := helper.NewFileRepo()
	if err := fileRepo.InitFromSource("/home/eiko/prog/git/fdf", dummy, dummy); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("files in repo: %d", fileRepo.Size())
	}
}
