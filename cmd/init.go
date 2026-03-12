package cmd

import (
	"fmt"

	"github.com/marco-souza/maestro/internal/scaffold"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new maestro project",
		Long:  "Set up a squad of OpenCode sub-agents, create maestro.yaml, and initialize the .maestro/ state directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing maestro project...")

			if err := scaffold.Init("."); err != nil {
				return err
			}

			fmt.Println("Done. Run 'maestro <objective>' to start a workflow.")
			return nil
		},
	}
}
