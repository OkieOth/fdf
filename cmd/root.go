package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sourceDir string
var searchRoot string
var whiteList []string
var blackList []string

var RootCmd = &cobra.Command{
	Use:   "fdf",
	Short: "find double files",
	Long: `The program takes a source directory and parses a destination directory
	for file duplicates found in source. There are several sub commands how to handle
	the found duplicates.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To do something useful, call the program with one of the available subcommands. To get a list of available sub-commands use the '--help' flag.")
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
	RootCmd.PersistentFlags().StringVar(&sourceDir, "source", "", "Source directory to start for finding required files")
	RootCmd.PersistentFlags().StringVar(&searchRoot, "searchRoot", "", "Root directory to search for duplicates")
	RootCmd.PersistentFlags().StringSliceVar(&whiteList, "whitelist", make([]string, 0), "Files to include via whitelist")
	RootCmd.PersistentFlags().StringSliceVar(&blackList, "blacklist", make([]string, 0), "Files to exclude via blacklist")
	RootCmd.MarkPersistentFlagDirname("source")
	RootCmd.MarkPersistentFlagDirname("searchRoot")
}

func arePersistentFlagsOk() (bool, []string) {
	checkDir := func(flagName string, dir string, isOk *bool, messages *[]string) {
		if dir == "" {
			*messages = append(*messages, fmt.Sprintf("'%s' flag is required, but is missing", flagName))
			*isOk = false
		} else {
			if fileInfo, err := os.Stat(dir); err != nil {
				*messages = append(*messages, fmt.Sprintf("can't access '%s' (%s): %v", flagName, dir, err))
				*isOk = false
			} else {
				if !fileInfo.IsDir() {
					*messages = append(*messages, fmt.Sprintf("'%s' (%s) seems to be no directory, but is required as one", flagName, sourceDir))
					*isOk = false
				}
			}
		}
	}
	messages := make([]string, 0)
	isOk := true
	checkDir("sourceDir", sourceDir, &isOk, &messages)
	checkDir("searchRoot", searchRoot, &isOk, &messages)
	return isOk, messages // TODO
}
