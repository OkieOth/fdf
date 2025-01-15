package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/okieoth/fdf/internal/pkg/cmdimpl"
	"github.com/spf13/cobra"
)

var json bool
var outputFile string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints the found duplicates",
	Long:  "Writes the paths of the found duplicates either to stdout or to a file",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for the greet command
		if isOk, messages := arePersistentFlagsOk(); !isOk {
			fmt.Println(messages)
			os.Exit(1)
		}
		if cpus != 0 {
			runtime.GOMAXPROCS(cpus)
		}
		cmdimpl.ListImpl(sourceDir, searchRoot, blackList, whiteList, noProgress, json, outputFile, true)
	},
}

func init() {
	ListCmd.PersistentFlags().BoolVar(&json, "json", false, "Prints the output in JSON instead of plain text")
	ListCmd.PersistentFlags().StringVar(&outputFile, "outputFile", "", "Optional filename to print the output in")
}
