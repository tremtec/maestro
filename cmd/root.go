package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "maestro",
		Short: "Maestro CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🎵 Welcome to Maestro")
			_ = cmd.Help()
		},
	}

	root.AddCommand(newInitCmd())

	return root
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
