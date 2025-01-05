package cmd

import (
	"fmt"
	"os"

	"github.com/okieoth/fdf/internal/pkg/cmdimpl"
	"github.com/spf13/cobra"
)

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
		cmdimpl.ListImpl(sourceDir, searchRoot, blackList, whiteList)
	},
}
