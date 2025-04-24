package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/samber/lo"
)

func updateWorkspacePackages(newPackages, oldPackages []string) error {
	useNew := lo.Map(newPackages, func(p string, _ int) string { return fmt.Sprintf("-use=%s", p) })
	dropOld := lo.Map(oldPackages, func(p string, _ int) string { return fmt.Sprintf("-dropuse=%s", p) })
	baseArgs := []string{"work", "edit"}
	fullArgs := append(baseArgs, dropOld...)
	fullArgs = append(fullArgs, useNew...)
	output, err := exec.Command("go", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	return nil
}

func updatePackageMods(c *PublishConfig, packages []string) error {
	for _, pkgn := range packages {
		workName := c.Packages[pkgn].WorkName
		pkgName := c.Packages[pkgn].PkgName
		// rename package name
		gomodPath := filepath.Join(pkgName, "go.mod")
		data, err := os.ReadFile(gomodPath)
		if err != nil {
			return fmt.Errorf("error reading go.mod file: %w", err)
		}
		data = bytes.Replace(data, []byte(workName), []byte(pkgName), 1)
		err = os.WriteFile(gomodPath, data, 0)
		if err != nil {
			return fmt.Errorf("error writing go.mod file: %w", err)
		}
	}
	return nil
}
