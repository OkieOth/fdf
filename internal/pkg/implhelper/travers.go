package implhelper

import (
	"fmt"
	"os"
	"path/filepath"
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

func TraversDir(dir string, blackList []string, whiteList []string, foundChan chan<- TraversResponse) {
	defer close(foundChan)
	handleFile := func(fileName string, wg *sync.WaitGroup, foundChan chan<- TraversResponse) {
		defer wg.Done()
		fmt.Println("handleFile: ", fileName) // FIXME, DEBUG
		// TODO - calc md5 and report it back
	}
	handleDir := func(dirName string, wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("filepath.Walk-1: ", dirName)
		err := filepath.WalkDir(dirName, func(path string, entry os.DirEntry, err error) error {
			//fmt.Println("filepath.Walk-2: ", path)
			if err != nil {
				return err
			} else {
				if path != dirName {
					if !entry.IsDir() {
						wg.Add(1)
						//fmt.Println("XXX write to handleFileChan: ", path)
						go handleFile(path, wg, foundChan)
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
