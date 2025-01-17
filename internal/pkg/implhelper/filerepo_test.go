package implhelper_test

import (
	"strings"
	"testing"

	"github.com/okieoth/fdf/internal/pkg/implhelper"
	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
)

func TestInitFromSource(t *testing.T) {
	implhelper.FileSizeThresholdInMB = 5
	whiteList := make([]string, 0)
	blackList := make([]string, 0)
	fileRepo := helper.NewFileRepo()
	fileRepo.InitFromSource("../..", blackList, whiteList, true)
	rs := fileRepo.Size()
	if rs == 0 {
		t.Errorf("Seems the fileRepo isn't initialized: %d", rs)
	}
}

func TestInitFromSourceBlackListed(t *testing.T) {
	implhelper.FileSizeThresholdInMB = 5
	whiteList := make([]string, 0)
	blackList := []string{"*.mod", "*.sum", "LICENSE", ".git", "README"}
	fileRepo := helper.NewFileRepo()
	fileRepo.InitFromSource("../../..", blackList, whiteList, true)
	rs := fileRepo.Size()
	if rs == 0 {
		t.Errorf("Seems the fileRepo isn't initialized: %d", rs)
	}
	readmeMd5, e := helper.GetMd5("../../../README.md")
	if e != nil {
		t.Errorf("Error while get md5 for README: %v", e)
	}
	if fileRepo.HasEntry(readmeMd5) {
		t.Error("seems init with blacklist isn't working 1")
	}
	mainMd5, e := helper.GetMd5("../../../main.go")
	if e != nil {
		t.Errorf("Error while get md5 for main.go: %v", e)
	}
	if !fileRepo.HasEntry(mainMd5) {
		t.Error("seems init with blacklist isn't working 2")
	}
}

func TestInitFromSourceWhiteListed(t *testing.T) {
	implhelper.FileSizeThresholdInMB = 5
	whiteList := []string{"*.go"}
	blackList := []string{}
	fileRepo := helper.NewFileRepo()
	fileRepo.InitFromSource("../../..", blackList, whiteList, true)
	rs := fileRepo.Size()
	if rs == 0 {
		t.Errorf("Seems the fileRepo isn't initialized: %d", rs)
	}
	for _, v := range fileRepo.Repo() {
		if !strings.HasSuffix(v.SourceFile, ".go") {
			t.Errorf("Found non-go file in the repo: %s", v.SourceFile)
		}
	}
}
