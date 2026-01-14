package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/github"
)

var (
	backupOutput string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of installed configurations",
	Long: `Create a timestamped backup archive of ~/.github/copilot configurations.

The backup is saved as a gzipped tar archive (.tar.gz) containing all
installed agents and skills.

Examples:
  cp-config backup                           Create backup in current directory
  cp-config backup --output ~/backups        Create backup in specific directory
  cp-config backup -o ./my-backup.tar.gz     Create backup with specific filename`,
	Run: runBackup,
}

func init() {
	backupCmd.Flags().StringVarP(&backupOutput, "output", "o", "", "Output path (directory or filename)")
	rootCmd.AddCommand(backupCmd)
}

func runBackup(cmd *cobra.Command, args []string) {
	// Get target directory
	target, err := github.NewTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if target exists
	if _, err := os.Stat(target.Root); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: no configurations found at %s\n", target.Root)
		os.Exit(1)
	}

	// Determine output path
	outputPath := determineBackupPath(backupOutput)

	// Create backup
	if err := createBackup(target.Root, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating backup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backup created: %s\n", outputPath)
}

func determineBackupPath(output string) string {
	timestamp := time.Now().Format("20060102-150405")
	defaultName := fmt.Sprintf("copilot-backup-%s.tar.gz", timestamp)

	if output == "" {
		return defaultName
	}

	// Check if output is a directory
	info, err := os.Stat(output)
	if err == nil && info.IsDir() {
		return filepath.Join(output, defaultName)
	}

	// If it ends with .tar.gz, use as-is
	if filepath.Ext(output) == ".gz" {
		return output
	}

	// Otherwise treat as directory and append default name
	return filepath.Join(output, defaultName)
}

func createBackup(sourceDir, outputPath string) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating backup file: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Walk the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If it's a file, write contents
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})
}
