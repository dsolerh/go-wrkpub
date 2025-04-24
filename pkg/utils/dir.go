package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dsolerh/go-wrkpub/pkg/config"
)

func copyDirectory(src, dst string) error {
	// The -r flag makes cp recursive
	// The -p flag preserves mode, ownership, and timestamps
	cmd := exec.Command("cp", "-rp", src, dst)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error copying from %s to %s: %v, output: %s", src, dst, err, output)
	}

	return nil
}

func CopyPackagesToRoot(c *config.PublishConfig, packages []string) error {
	for _, pkg := range packages {
		workPath := filepath.Join(c.Root, c.Packages[pkg].WorkName)
		pkgPath := filepath.Join(c.Root, c.Packages[pkg].PkgName)
		if err := copyDirectory(workPath, pkgPath); err != nil {
			return err
		}
	}
	return nil
}

func RemovePackagesFromRoot(packagesName []string) error {
	for _, packageName := range packagesName {
		// remove the package from the root
		if err := os.RemoveAll(packageName); err != nil {
			return fmt.Errorf("error removing package directory: %w\n", err)
		}
	}
	return nil
}
