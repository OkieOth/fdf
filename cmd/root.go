package cmd

import (
	"fmt"
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
	RootCmd.MarkPersistentFlagRequired("source")
	RootCmd.MarkPersistentFlagDirname("source")
	RootCmd.MarkPersistentFlagRequired("searchRoot")
	RootCmd.MarkPersistentFlagDirname("searchRoot")
}

func persistentFlagsAreOk() bool {
	return false // TODO
}
