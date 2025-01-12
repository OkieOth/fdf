package cmdimpl_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/okieoth/fdf/internal/pkg/cmdimpl"
)

func TestListImpl_1(t *testing.T) {
	whiteList := make([]string, 0)
	blackList := make([]string, 0)
	sourceDir := "../.."
	searchRoot := "../../.."

	if fileRepo, err := cmdimpl.ListImpl(sourceDir, searchRoot, blackList, whiteList, true); err != nil {
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

func TestListImpl_2(t *testing.T) {
	whiteList := []string{"list_test.go"}
	blackList := make([]string, 0)
	sourceDir := "../.."
	searchRoot := "../../.."

	if fileRepo, err := cmdimpl.ListImpl(sourceDir, searchRoot, blackList, whiteList, true); err != nil {
		t.Errorf("Error while searching for duplicates: %v", err)
	} else {
		if fileRepo.Size() != 1 {
			t.Errorf("Wrong number of fileRepo entries. Expected: 1, got %d", fileRepo.Size())
		} else {
			repo := fileRepo.Repo()
			count := 0
			for _, fre := range repo {
				count++
				if len(fre.Duplicates) != 1 {
					t.Errorf("No duplicates for file: %s", fre.SourceFile)
				}
				src := fre.SourceFile[5:]
				if !strings.HasSuffix(fre.Duplicates[0], src) {
					t.Errorf("Wrong duplicate: src=%s, found=%s", fre.SourceFile, fre.Duplicates[0])
				}
			}
			if count != 1 {
				t.Error("something went wrong")
			}
		}
	}
}

func TestListImpl_3(t *testing.T) {
	whiteList := make([]string, 0)
	blackList := []string{".git"}
	sourceDir := "../../.."

	if fileRepo, err := cmdimpl.ListImpl(sourceDir, "", blackList, whiteList, true); err != nil {
		t.Errorf("Error while searching for duplicates: %v", err)
	} else {
		fmt.Println("FileRepo.Size: ", fileRepo.Size())
		repo := fileRepo.Repo()
		for _, fre := range repo {
			if fre.SourceFile == "../../../README.md" {
				if len(fre.Duplicates) != 1 {
					t.Errorf("No duplicates for file: %s", fre.SourceFile)
				}
				if fre.Duplicates[0] != "../../../docs/README.md" {
					t.Errorf("Wrong duplicate: %s", fre.Duplicates[0])
				}
			}
		}
	}
}
