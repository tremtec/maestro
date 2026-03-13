package cmd

import (
	"fmt"
	"strings"

	"github.com/marco-souza/maestro/internal/scaffold"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var tools []string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new maestro project",
		Long: fmt.Sprintf(
			"Set up a squad of sub-agents, create maestro.yaml, and initialize the .maestro/ state directory.\n\n"+
				"Supported tools: %s",
			strings.Join(scaffold.SupportedTools(), ", "),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing maestro project...")

			if err := scaffold.Init(".", tools...); err != nil {
				return err
			}

			fmt.Println("Done. Run 'maestro <objective>' to start a workflow.")
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tools, "tool", []string{"opencode"},
		"agent runtimes to scaffold (e.g., --tool opencode,amp)")

	return cmd
}
