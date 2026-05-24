package ln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "ln",
	Usage: "Make links between files",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "symbolic", Aliases: []string{"s"}, Usage: "Create symbolic link instead of hard link"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		symbolic := cmd.Bool("symbolic")

		files := cmd.Args().Slice()
		if len(files) < 2 {
			return fmt.Errorf("missing file operand")
		}

		target, _ := filepath.Abs(files[0])
		linkName, _ := filepath.Abs(files[1])

		if symbolic {
			return createSymlink(target, linkName)
		}
		return createHardlink(target, linkName)
	},
}

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
	if _, err := os.Stat(target); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("target '%s' does not exist", target)
		}
		return fmt.Errorf("cannot access target '%s': %v", target, err)
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
			"failed to create hard link '%s' -> '%s': %v (target may not be on NTFS)",
			linkName, target, err,
		)
	}
	return nil
}
