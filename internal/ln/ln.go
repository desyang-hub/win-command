package ln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "ln",
	Usage: "Make links between files",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "symbolic, s", Usage: "Create symbolic link instead of hard link"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		symbolic := cmd.Bool("symbolic")

		files := cmd.Args().Slice()
		if len(files) < 2 {
			return fmt.Errorf("missing file operand")
		}

		target := files[0]
		linkName := files[1]

		target, _ = filepath.Abs(target)

		if runtime.GOOS == "windows" {
			if symbolic {
				return createSymlink(target, linkName)
			}
			return createHardlink(target, linkName)
		}

		if symbolic {
			return os.Symlink(target, linkName)
		}
		return os.Link(target, linkName)
	},
}

func createSymlink(target, linkName string) error {
	target, _ = filepath.Abs(target)
	os.Remove(linkName)

	var err error
	if runtime.GOOS == "windows" {
		err = os.Symlink(target, linkName)
	} else {
		err = os.Symlink(target, linkName)
	}
	if err != nil {
		return fmt.Errorf("failed to create symlink: %v (may need admin or Developer Mode on Windows)", err)
	}
	return nil
}

func createHardlink(target, linkName string) error {
	target, _ = filepath.Abs(target)
	if err := os.Link(target, linkName); err != nil {
		return fmt.Errorf("failed to create hard link: %v (target may not be on NTFS)", err)
	}
	return nil
}
