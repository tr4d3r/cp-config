package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/github"
)

var (
	restoreForce bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore <backup-file>",
	Short: "Restore configurations from a backup",
	Long: `Restore ~/.github/copilot configurations from a backup archive.

This command extracts the contents of a backup archive created with
'cp-config backup' to the Copilot configuration directory.

Use --force to overwrite existing configurations without prompting.

Examples:
  cp-config restore copilot-backup-20240115-120000.tar.gz
  cp-config restore ~/backups/my-backup.tar.gz --force`,
	Args: cobra.ExactArgs(1),
	Run:  runRestore,
}

func init() {
	restoreCmd.Flags().BoolVarP(&restoreForce, "force", "f", false, "Overwrite existing configurations without prompting")
	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) {
	backupPath := args[0]

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: backup file not found: %s\n", backupPath)
		os.Exit(1)
	}

	// Get target directory
	target, err := github.NewTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if target already has configurations
	if !restoreForce {
		hasConfigs, err := targetHasConfigs(target.Root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking target: %v\n", err)
			os.Exit(1)
		}

		if hasConfigs {
			fmt.Printf("Warning: configurations already exist at %s\n", target.Root)
			fmt.Print("Overwrite? [y/N] ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" {
				fmt.Println("Aborted.")
				return
			}
		}
	}

	// Restore backup
	if err := restoreBackup(backupPath, target.Root); err != nil {
		fmt.Fprintf(os.Stderr, "Error restoring backup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Restored configurations to %s\n", target.Root)
}

func targetHasConfigs(targetDir string) (bool, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return false, nil
	}

	// Check for any .md files in agents or skills directories
	agentsDir := filepath.Join(targetDir, "agents")
	skillsDir := filepath.Join(targetDir, "skills")

	for _, dir := range []string{agentsDir, skillsDir} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return false, err
		}

		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
				return true, nil
			}
		}
	}

	return false, nil
}

func restoreBackup(backupPath, targetDir string) error {
	// Open backup file
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("opening backup file: %w", err)
	}
	defer file.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("reading gzip: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}

		// Construct target path
		targetPath := filepath.Join(targetDir, header.Name)

		// Ensure we don't escape the target directory (security check)
		if !isSubPath(targetDir, targetPath) {
			return fmt.Errorf("invalid path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}

		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("creating parent directory: %w", err)
			}

			// Create file
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("creating file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("writing file: %w", err)
			}
			outFile.Close()
		}
	}

	return nil
}

// isSubPath checks if child is under parent (prevents path traversal)
func isSubPath(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	// Check that we don't escape parent with ..
	return rel != ".." && !filepath.IsAbs(rel) && len(rel) >= 1 && rel[0] != '.'
}
