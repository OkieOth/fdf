package cmdimpl

import (
	"fmt"

	"github.com/okieoth/fdf/internal/pkg/implhelper"
	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
	"github.com/okieoth/fdf/internal/pkg/progressbar"
)

func ListImpl(sourceDir string, searchRoot string, blackList []string, whiteList []string, noProgress bool) (*helper.FileRepo, error) {
	fileRepo := helper.NewFileRepo()
	if noProgress {
		progressbar.Init(0, "fdf ist running")
	} else {
		progressbar.Init(implhelper.GetFileCount(sourceDir, searchRoot, blackList, whiteList), "Init file repo")
	}
	if err := fileRepo.InitFromSource(sourceDir, blackList, whiteList, noProgress); err != nil {
		return fileRepo, fmt.Errorf("Error while initialize from sourceDir: %v\n", err)
	}
	if searchRoot != "" {
		doneChan := make(chan *error)
		if !noProgress {
			progressbar.Description("Search for duplicates")
		}
		go implhelper.SearchForDuplicates(searchRoot, blackList, whiteList, fileRepo, doneChan, noProgress)
		for e := range doneChan {
			if e != nil {
				return fileRepo, fmt.Errorf("Error while search for duplicates: %v", e)
			}
		}
	}
	if fileRepo.HasDuplicates() {
		fmt.Println("Found file duplicates")
		fmt.Println("=================================")
		for _, fre := range fileRepo.Repo() {
			if len(fre.Duplicates) > 0 {
				fmt.Println("File: ", fre.SourceFile)
				for _, d := range fre.Duplicates {
					fmt.Println("    ", d)
				}
				fmt.Println("---------------------------------")
			}
		}
	} else {
		fmt.Println("No duplicates found")
	}
	fmt.Println()
	return fileRepo, nil
}
