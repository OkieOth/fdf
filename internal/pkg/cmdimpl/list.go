package cmdimpl

import (
	"fmt"

	"github.com/okieoth/fdf/internal/pkg/implhelper"
	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
)

func ListImpl(sourceDir string, searchRoot string, blackList []string, whiteList []string) (*helper.FileRepo, error) {
	fileRepo := helper.NewFileRepo()
	if err := fileRepo.InitFromSource(sourceDir, blackList, whiteList); err != nil {
		return fileRepo, fmt.Errorf("Error while initialize from sourceDir: %v\n", err)
	}
	if searchRoot != "" {
		doneChan := make(chan *error)
		go implhelper.SearchForDuplicates(searchRoot, blackList, whiteList, fileRepo, doneChan)
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
