package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints the found duplicates",
	Long:  "Writes the paths of the found duplicates either to stdout or to a file",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for the greet command
		fmt.Println("TODO - List found duplicate files")
	},
}
