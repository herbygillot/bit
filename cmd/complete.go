package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Add classical tab completion to bit",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		installClassicalTabCompletion()
	},
	Args: cobra.NoArgs,
}

func init() {
	ShellCmd.AddCommand(completeCmd)
}

func installClassicalTabCompletion() {
	os.Setenv("COMP_INSTALL", "1")
	Bitcomplete()
}
