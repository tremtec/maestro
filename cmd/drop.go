package cmd

import (
	"fmt"

	"github.com/marco-souza/maestro/internal/scaffold"
	"github.com/spf13/cobra"
)

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Unset Maestro configuration values",
	Long: `Drop allows you to unset Maestro configuration values. For example, you can use it to remove a tool from your setup:
	maestro drop
This will remove the opencode tool from your Maestro configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Dropping maestro set up...")

		if err := scaffold.Drop("."); err != nil {
			return err
		}

		fmt.Println("Done.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dropCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dropCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
