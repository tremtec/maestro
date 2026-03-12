package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new maestro project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing maestro project...")
			fmt.Println("Done.")
			return nil
		},
	}
}
