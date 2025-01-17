package cmdimpl

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/okieoth/fdf/internal/pkg/implhelper"
	helper "github.com/okieoth/fdf/internal/pkg/implhelper"
	"github.com/okieoth/fdf/internal/pkg/progressbar"
)

func ListImpl(sourceDir string, searchRoot string, blackList []string, whiteList []string, noProgress bool, json bool, outputFile string, ignoreSameFiles bool) (*helper.FileRepo, error) {
	fileRepo := helper.NewFileRepo()
	if noProgress {
		progressbar.Init(0, "fdf ist running")
	} else {
		progressbar.Init(implhelper.GetFileCount(sourceDir, searchRoot, blackList, whiteList), "Init file repo")
	}
	fileRepo.InitFromSource(sourceDir, blackList, whiteList, noProgress)
	if searchRoot != "" {
		doneChan := make(chan *error)
		if !noProgress {
			progressbar.Description("Search for duplicates")
		}
		go implhelper.SearchForDuplicates(searchRoot, blackList, whiteList, fileRepo, doneChan, noProgress, ignoreSameFiles)
		for e := range doneChan {
			if e != nil {
				return fileRepo, fmt.Errorf("error while search for duplicates: %v", e)
			}
		}
	}
	return printOutput(fileRepo, json, outputFile)
}

func extractDuplicates(repo map[string]helper.FileRepoEntry) *[]helper.FileRepoEntry {
	ret := make([]helper.FileRepoEntry, 0)
	for _, fre := range repo {
		if len(fre.Duplicates) > 0 {
			ret = append(ret, fre)
		}
	}
	return &ret
}

func printOutput(fileRepo *helper.FileRepo, jsonOutput bool, outputFilePath string) (*helper.FileRepo, error) {
	if fileRepo.HasDuplicates() {
		if jsonOutput {
			duplicates := extractDuplicates(fileRepo.Repo())
			jsonData, err := json.MarshalIndent(duplicates, "", "  ")
			if err != nil {
				return fileRepo, fmt.Errorf("error while marshalling fileRepo: %v", err)
			}
			if outputFilePath != "" {
				outputFile, err := os.Create(outputFilePath)
				if err != nil {
					return fileRepo, fmt.Errorf("error while creating output file: %v", err)
				}
				defer outputFile.Close()

				_, err = outputFile.Write(jsonData)
				if err != nil {
					return fileRepo, fmt.Errorf("error while writing output file: %v", err)
				}
			} else {
				fmt.Println(jsonData)
			}
		} else {
			if outputFilePath != "" {
				outputFile, err := os.Create(outputFilePath)
				if err != nil {
					return fileRepo, fmt.Errorf("error while creating output file: %v", err)
				}
				defer outputFile.Close()
				outputFile.WriteString("Found file duplicates\n")
				outputFile.WriteString("=================================\n")
				for _, fre := range fileRepo.Repo() {
					if len(fre.Duplicates) > 0 {
						outputFile.WriteString(fmt.Sprintf("File: %s\n", fre.SourceFile))
						for _, d := range fre.Duplicates {
							outputFile.WriteString(fmt.Sprintf("    %s\n", d))
						}
						outputFile.WriteString("---------------------------------\n")
					}
				}
			} else {
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
				fmt.Println()
			}
		}
	} else {
		fmt.Println("No duplicates found")
		fmt.Println()
	}
	return fileRepo, nil
}
