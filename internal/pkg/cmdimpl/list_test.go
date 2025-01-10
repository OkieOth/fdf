package cmdimpl_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/okieoth/fdf/internal/pkg/cmdimpl"
)

func TestListImpl(t *testing.T) {
	whiteList := make([]string, 0)
	blackList := make([]string, 0)
	sourceDir := "../.."
	searchRoot := "../../.."

	if fileRepo, err := cmdimpl.ListImpl(sourceDir, searchRoot, blackList, whiteList); err != nil {
		t.Errorf("Error while searching for duplicates: %v", err)
	} else {
		fmt.Println("FileRepo.Size: ", fileRepo.Size())
		repo := fileRepo.Repo()
		for _, fre := range repo {
			if len(fre.Duplicates) != 1 {
				t.Errorf("No duplicates for file: %s", fre.SourceFile)
			}
			src := fre.SourceFile[5:]
			if !strings.HasSuffix(fre.Duplicates[0], src) {
				t.Errorf("Wrong duplicate: src=%s, found=%s", fre.SourceFile, fre.Duplicates[0])
			}
		}
	}
}
