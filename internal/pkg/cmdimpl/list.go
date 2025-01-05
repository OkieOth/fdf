package cmdimpl

import (
	"fmt"

	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
)

func ListImpl(sourceDir string, searchRoot string, blackList []string, whiteList []string) {
	fileRepo := helper.NewFileRepo()
	if err := fileRepo.InitFromSource(sourceDir, blackList, whiteList); err != nil {
		fmt.Printf("Error while initialize from sourceDir: %v\n", err)
		return
	}
	// TODO
}
