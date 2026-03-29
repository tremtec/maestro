package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tremtec/maestro/internal/scaffold"
)

var tools []string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new maestro project",
	Long: `Set up a squad of sub-agents, create maestro.yaml, and initialize the .maestro/ state directory.

Supported tool: opencode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing maestro project...")

		if err := scaffold.Init(".", tools...); err != nil {
			return err
		}

		fmt.Println("Done. Run 'maestro run \"<prompt>\"' to start a workflow.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringSliceVarP(&tools, "tool", "t", []string{"opencode"}, "agent runtime to scaffold (only opencode is supported)")
}
