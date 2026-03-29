package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the maestro CLI to the latest version",
	Long: `Download and install the latest version of maestro from GitHub releases.

This command will:
- Check for the latest release on GitHub
- Download the appropriate binary for your platform
- Replace the current binary with the new version

Note: You may need sudo privileges if maestro is installed in a system directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for updates...")

		// Get current version
		currentVersion, err := getCurrentVersion()
		if err != nil {
			return fmt.Errorf("getting current version: %w", err)
		}

		// Get latest version from GitHub
		latestVersion, downloadURL, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("checking for updates: %w", err)
		}

		// Compare versions
		if currentVersion == latestVersion {
			fmt.Printf("You're already on the latest version: v%s\n", currentVersion)
			return nil
		}

		fmt.Printf("Upgrading from v%s to v%s...\n", currentVersion, latestVersion)

		// Download the new binary
		if err := downloadAndInstall(downloadURL); err != nil {
			return fmt.Errorf("upgrading: %w", err)
		}

		fmt.Println("Upgrade complete!")
		return nil
	},
}

func getCurrentVersion() (string, error) {
	// Try to get version from the binary itself
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	cmd := exec.CommandContext(ctx, "maestro", "version")
	output, err := cmd.Output()
	if err != nil {
		// If version command doesn't exist, return empty
		return "0.0.0", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func getLatestRelease() (version, downloadURL string, err error) {
	// Fetch latest release info from GitHub
	resp, err := http.Get("https://api.github.com/repos/tremtec/maestro/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Simple parsing - extract tag_name and find binary URL
	content := string(body)

	// Extract version from tag_name
	tagIdx := strings.Index(content, `"tag_name":"`)
	if tagIdx == -1 {
		return "", "", fmt.Errorf("could not find tag_name in response")
	}
	tagStart := tagIdx + len(`"tag_name":"`)
	tagEnd := strings.Index(content[tagStart:], `"`)
	version = content[tagStart : tagStart+tagEnd]
	version = strings.TrimPrefix(version, "v")

	// Determine binary name for current platform
	goos := runtime.GOOS
	arch := runtime.GOARCH

	var binaryName string
	switch goos {
	case "darwin":
		if arch == "arm64" {
			binaryName = "maestro-darwin-arm64"
		} else {
			binaryName = "maestro-darwin-amd64"
		}
	case "linux":
		if arch == "arm64" {
			binaryName = "maestro-linux-arm64"
		} else if arch == "386" {
			binaryName = "maestro-linux-386"
		} else {
			binaryName = "maestro-linux-amd64"
		}
	case "windows":
		if arch == "386" {
			binaryName = "maestro-windows-386.exe"
		} else {
			binaryName = "maestro-windows-amd64.exe"
		}
	default:
		return "", "", fmt.Errorf("unsupported platform: %s/%s", goos, arch)
	}

	// Find the asset URL
	downloadURL = fmt.Sprintf("https://github.com/tremtec/maestro/releases/download/v%s/%s.tar.gz", version, binaryName)

	return version, downloadURL, nil
}

func downloadAndInstall(downloadURL string) error {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "maestro-upgrade-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Download
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("saving download: %w", err)
	}

	// Get current binary path
	currentPath, err := exec.LookPath("maestro")
	if err != nil {
		// If not in PATH, use the directory we're running from
		currentPath, _ = os.Executable()
	}

	// Replace the binary (note: this requires appropriate permissions)
	if currentPath != "" {
		// Backup current binary
		backupPath := currentPath + ".bak"
		if err := os.Rename(currentPath, backupPath); err != nil {
			return fmt.Errorf("backing up current binary: %w", err)
		}

		// Move new binary to location
		newBinaryPath := filepath.Join(filepath.Dir(currentPath), "maestro")
		if err := os.Rename(tmpFile.Name(), newBinaryPath); err != nil {
			// Restore backup - ignore error since we're already in error state
			_ = os.Rename(backupPath, currentPath)
			return fmt.Errorf("installing new binary: %w", err)
		}

		// Make executable
		if err := os.Chmod(newBinaryPath, 0o755); err != nil {
			return fmt.Errorf("chmod new binary: %w", err)
		}

		// Remove backup - ignore error
		_ = os.Remove(backupPath)
	} else {
		return fmt.Errorf("could not find maestro binary location")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
