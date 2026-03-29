package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tremtec/maestro/internal/scaffold"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update prompt agents and file structure",
	Long: `Update the agent definitions and file structure to the latest version.

This command will:
- Re-scaffold all agent files from the latest templates
- Update the .opencode/ directory with current agent prompts
- Preserve existing customizations in .maestro/

Run this command when you want to sync with the latest Maestro agent definitions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Updating maestro agents and file structure...")

		targetDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		// Check if project is initialized
		if !scaffold.IsInitialized(targetDir) {
			fmt.Println("Project not initialized. Run 'maestro init' first.")
			return nil
		}

		// Re-scaffold opencode (the only supported tool now)
		if err := scaffold.Update(targetDir, "opencode"); err != nil {
			return err
		}

		fmt.Println("Update complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
