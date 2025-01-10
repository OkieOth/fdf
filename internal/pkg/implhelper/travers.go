package implhelper

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

func TraversDir(dir string, blackList []string, whiteList []string, foundChan chan<- TraversResponse) {
	defer close(foundChan)
	handleFile := func(fileName string, wg *sync.WaitGroup, foundChan chan<- TraversResponse) {
		defer wg.Done()
		if md5, err := GetMd5(fileName); err == nil {
			foundChan <- FoundTraversResponse(fileName, md5)
		} else {
			foundChan <- ErrorTraversResponse(err)
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

func SearchForDuplicates(searchRoot string, blackList []string, whiteList []string, fileRepo *FileRepo, doneChan chan<- *error) {
	defer close(doneChan)
	resp := make(chan TraversResponse)
	go TraversDir(searchRoot, blackList, whiteList, resp)
	for r := range resp {
		if r.err != nil {
			doneChan <- &r.err
			return
		} else {
			fileRepo.CheckForDuplicateAndAddInCase(r.md5, r.file)
		}
	}
}
