package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// migrateUpCmd represents the migrateUp command
var migrateUpCmd = &cobra.Command{
	Use:   "migrateUp",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateUp called")
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)
}
