package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// HasUncommittedChanges checks if the git repository at the given path
// has any uncommitted changes.
func HasUncommittedChanges() (bool, error) {
	// Run: git status --porcelain
	output, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, fmt.Errorf("error running git status: %w", err)
	}

	// If the output is not empty, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func CommitChanges(message string) error {
	// add the changes
	output, err := exec.Command("git", "add", ".").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to git: %w, output: %s\n", err, output)
	}
	output, err = exec.Command("git", "commit", "-m", message).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error commiting the changes: %w, output: %s\n", err, output)
	}
	return nil
}

func GetPublishCommitMessage(pkgtags []string) string {
	var plural = "package"
	if len(pkgtags) > 1 {
		plural = "packages"
	}
	return fmt.Sprintf("ci: publish %s\n +%s", plural, strings.Join(pkgtags, "\n +"))
}

func TagPackagesVersion(tags []string) error {
	for _, tag := range tags {
		output, err := exec.Command("git", "tag", "--", tag, "HEAD").CombinedOutput()
		if err != nil {
			return fmt.Errorf("error tagging the changes: %w, output: %s\n", err, output)
		}
	}

	return nil
}

const CleanupCommit = "ci: cleanup publish"

func PushChanges(tags []string) error {
	// push the changes
	output, err := exec.Command("git", "push").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the changes: %w, output: %s\n", err, output)
	}

	baseArgs := []string{"push", "origin"}
	fullArgs := append(baseArgs, tags...)
	output, err = exec.Command("git", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging the changes: %w, output: %s\n", err, output)
	}
	return nil
}
