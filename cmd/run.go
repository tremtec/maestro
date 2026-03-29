package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	opencodeURL string
	sessionID   string
)

var runCmd = &cobra.Command{
	Use:   "run \"<prompt>\"",
	Short: "Run a maestro session on the opencode server",
	Long: `Send a prompt to the opencode server to execute a maestro session.

This command communicates with the opencode API server to create and run
an agentic workflow. By default, it connects to http://localhost:8080.

You can customize the server URL using the --url flag.

Example:
  maestro run "Create a new API endpoint for user authentication"
  maestro run --url https://opencode.example.com "Build a feature login"
  maestro run --session abc123 "Continue the previous session"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := args[0]

		// Check if opencode CLI is available (as alternative to HTTP API)
		if isOpencodeCLIAvailable() {
			return runWithOpencodeCLI(prompt)
		}

		// Otherwise use HTTP API
		return runWithHTTPAPI(prompt)
	},
}

func isOpencodeCLIAvailable() bool {
	_, err := exec.LookPath("opencode")
	return err == nil
}

func runWithOpencodeCLI(prompt string) error {
	fmt.Println("Using opencode CLI...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Run opencode with the prompt
	opencodeCmd := exec.CommandContext(ctx, "opencode", prompt)
	opencodeCmd.Stdout = os.Stdout
	opencodeCmd.Stderr = os.Stderr
	opencodeCmd.Stdin = os.Stdin

	return opencodeCmd.Run()
}

func runWithHTTPAPI(prompt string) error {
	// Default URL
	if opencodeURL == "" {
		opencodeURL = "http://localhost:8080"
	}

	// Build the request
	endpoint := opencodeURL + "/api/sessions"

	payload := map[string]interface{}{
		"prompt": prompt,
		"mode":   "maestro",
	}

	if sessionID != "" {
		payload["session_id"] = sessionID
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	// Make the request
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("connecting to opencode server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status: %d - %s", resp.StatusCode, string(body))
	}

	// Read and display response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	// Display session info
	if sid, ok := result["session_id"].(string); ok {
		fmt.Printf("Session ID: %s\n", sid)
	}

	if output, ok := result["output"].(string); ok {
		fmt.Println(output)
	}

	// Display any errors
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Errors:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %v\n", e)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&opencodeURL, "url", "", "Opencode server URL (default: http://localhost:8080)")
	runCmd.Flags().StringVar(&sessionID, "session", "", "Continue an existing session")
}
