package ln

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// Cmd is the ln command, initialized by platform-specific init()
var Cmd *cli.Command

func createSymlink(target, linkName string) error {
	os.Remove(linkName)
	if err := os.Symlink(target, linkName); err != nil {
		return fmt.Errorf(
			"failed to create symlink '%s' -> '%s': %v (may need admin or Developer Mode on Windows)",
			linkName, target, err,
		)
	}
	return nil
}

func createHardlink(target, linkName string) error {
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("target '%s' does not exist", target)
		}
		return fmt.Errorf("cannot access target '%s': %v", target, err)
	}

	if info.IsDir() {
		return fmt.Errorf("cannot create hard link to directory '%s' (use -s for symlink)", target)
	}

	parent := filepath.Dir(linkName)
	if _, err := os.Stat(parent); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("parent directory '%s' of link does not exist", parent)
		}
		return fmt.Errorf("cannot access parent directory '%s': %v", parent, err)
	}

	os.Remove(linkName)

	if err := os.Link(target, linkName); err != nil {
		return fmt.Errorf(
			"failed to create hard link '%s' -> '%s': %v",
			linkName, target, err,
		)
	}
	return nil
}
