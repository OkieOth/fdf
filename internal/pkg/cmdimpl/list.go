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
	doneChan := make(chan *error)
	go implhelper.SearchForDuplicates(searchRoot, blackList, whiteList, fileRepo, doneChan)
	for e := range doneChan {
		if e != nil {
			return fileRepo, fmt.Errorf("Error while search for duplicates: %v", e)
		}
	}
	return fileRepo, nil
}
