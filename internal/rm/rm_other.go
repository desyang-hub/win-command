//go:build !windows

package rm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "rm",
	Usage: "Remove files or directories",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "recursive, r", Usage: "Remove directories recursively"},
		&cli.BoolFlag{Name: "force, f", Usage: "Ignore nonexistent files, never prompt"},
		&cli.BoolFlag{Name: "permanent, P", Usage: "Permanently delete (skip recycle bin)"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		recursive := cmd.Bool("recursive")
		force := cmd.Bool("force")

		files := cmd.Args().Slice()
		if len(files) == 0 {
			return fmt.Errorf("missing file operand")
		}

		for _, file := range files {
			if err := remove(file, recursive, force); err != nil {
				return err
			}
		}
		return nil
	},
}

func remove(file string, recursive, force bool) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && force {
			return nil
		}
		return err
	}

	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("cannot remove '%s': it is a directory (use -r)", path)
		}
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}
