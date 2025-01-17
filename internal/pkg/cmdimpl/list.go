package cmdimpl

import (
	"encoding/json"
	"fmt"

	"golang.org/x/term"

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

func extentOutputLine(s string, termWith int) string {
	runeSlice := []rune(s)
	lastRune := runeSlice[len(runeSlice)-1]
	for i := len(runeSlice); i < termWith; i++ {
		s += string(lastRune)
	}
	return s
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
			const LINE_BASE_1 = "╔═══════ Source File ════"
			const LINE_BASE_2 = "║ "
			const LINE_BASE_3 = "╠─────── Duplicates ─────"
			const LINE_BASE_4 = "╚════════════════════════"
			const LINE_BASE_LEN = 25
			repo := fileRepo.Repo()
			maxPathPrefixLen := implhelper.GetMaxPathPrefixLen(repo)

			if outputFilePath != "" {
				outputFile, err := os.Create(outputFilePath)
				if err != nil {
					return fileRepo, fmt.Errorf("error while creating output file: %v", err)
				}
				defer outputFile.Close()
				width := implhelper.GetMaxTextLen(repo)
				if width < LINE_BASE_LEN+2 {
					width = LINE_BASE_LEN
				}
				line1 := extentOutputLine(LINE_BASE_1, width)
				line3 := extentOutputLine(LINE_BASE_3, width)
				line4 := extentOutputLine(LINE_BASE_4, width)

				for _, fre := range repo {
					if len(fre.Duplicates) > 0 {
						outputFile.WriteString(line1)
						outputFile.WriteString("\n")
						outputFile.WriteString(LINE_BASE_2)
						outputFile.WriteString(fre.SourceFile[maxPathPrefixLen:])
						outputFile.WriteString("\n")
						outputFile.WriteString(line3)
						outputFile.WriteString("\n")
						for _, d := range fre.Duplicates {
							outputFile.WriteString(LINE_BASE_2)
							outputFile.WriteString(d[maxPathPrefixLen:])
							outputFile.WriteString("\n")
						}
						outputFile.WriteString(line4)
						outputFile.WriteString("\n\n")
					}
				}
			} else {
				termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
				if err != nil || termWidth == 0 {
					termWidth = 100
				}
				line1 := extentOutputLine(LINE_BASE_1, termWidth)
				line3 := extentOutputLine(LINE_BASE_3, termWidth)
				line4 := extentOutputLine(LINE_BASE_4, termWidth)
				for _, fre := range repo {
					if len(fre.Duplicates) > 0 {
						fmt.Println(line1)
						fmt.Println(LINE_BASE_2, fre.SourceFile[maxPathPrefixLen:])
						fmt.Println(line3)
						for _, d := range fre.Duplicates {
							fmt.Println(LINE_BASE_2, d[maxPathPrefixLen:])
						}
						fmt.Println(line4)
						fmt.Println()
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
