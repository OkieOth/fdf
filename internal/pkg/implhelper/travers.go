package implhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/okieoth/fdf/internal/pkg/progressbar"
)

type TraversResponse struct {
	file string
	md5  string
	err  error
}

func FoundTraversResponse(file string, md5 string) TraversResponse {
	return TraversResponse{
		file: file,
		md5:  md5,
	}
}

func ErrorTraversResponse(err error) TraversResponse {
	return TraversResponse{
		err: err,
	}
}

func blackWhileListMatch(patternToMatch string, path string) bool {
	if strings.Contains(patternToMatch, "*") {
		// *.jpeg; IMG*.j*
		tmp := regexp.QuoteMeta(patternToMatch)
		regexpStr := strings.ReplaceAll(tmp, `\*`, ".*")
		r := regexp.MustCompile(regexpStr)
		return r.Match([]byte(path))

	} else {
		return strings.Contains(path, patternToMatch)
	}
}

func shouldBeProcessed(path string, blackList []string, whiteList []string) bool {
	if len(blackList) > 0 {
		for _, b := range blackList {
			if blackWhileListMatch(b, path) {
				return false
			}
		}
		return true
	} else {
		if len(whiteList) > 0 {
			for _, b := range whiteList {
				if blackWhileListMatch(b, path) {
					return true
				}
			}
			return false
		} else {
			return true
		}
	}
}

func TraversDir(dir string, blackList []string, whiteList []string, foundChan chan<- TraversResponse, skipMd5 bool, ignoreSameFiles bool) {
	defer close(foundChan)
	handleFile := func(fileName string, wg *sync.WaitGroup, foundChan chan<- TraversResponse) {
		defer wg.Done()
		if skipMd5 {
			foundChan <- FoundTraversResponse(fileName, "")
		} else {
			if md5, err := GetMd5(fileName); err == nil {
				foundChan <- FoundTraversResponse(fileName, md5)
			} else {
				foundChan <- ErrorTraversResponse(err)
			}
		}
	}
	handleDir := func(dirName string, wg *sync.WaitGroup) {
		defer wg.Done()
		err := filepath.WalkDir(dirName, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			} else {
				if path != dirName {

					if !entry.IsDir() {
						if shouldBeProcessed(path, blackList, whiteList) {
							wg.Add(1)
							if ignoreSameFiles {
								if absPath, err := filepath.Abs(path); err == nil {
									path = absPath
								} else {
									foundChan <- ErrorTraversResponse(fmt.Errorf("Error while getting fileInfo (path: %s): %v", path, err))
								}
							}
							go handleFile(path, wg, foundChan)
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			foundChan <- ErrorTraversResponse(err)
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go handleDir(dir, &wg)
	wg.Wait()
}

func SearchForDuplicates(searchRoot string, blackList []string, whiteList []string, fileRepo *FileRepo, doneChan chan<- *error, noProgress bool, ignoreSameFiles bool) {
	defer close(doneChan)
	resp := make(chan TraversResponse)
	empty := make([]string, 0)
	go TraversDir(searchRoot, empty, empty, resp, false, ignoreSameFiles)
	for r := range resp {
		if r.err != nil {
			doneChan <- &r.err
			return
		} else {
			fileRepo.CheckForDuplicateAndAddInCase(r.md5, r.file)
		}
		progressbar.ProgressOne()
	}
}

func GetFileCount(sourceDir string, searchRoot string, blackList []string, whiteList []string) int64 {
	ret := int64(0)
	traversDirAndCount := func(dir string, blackList []string, whiteList []string) {
		resp := make(chan TraversResponse)
		go TraversDir(dir, blackList, whiteList, resp, true, false)
		for _ = range resp {
			ret++
		}
	}
	traversDirAndCount(sourceDir, blackList, whiteList)
	if searchRoot != "" {
		empty := make([]string, 0)
		traversDirAndCount(searchRoot, empty, empty)
	}
	return ret
}
