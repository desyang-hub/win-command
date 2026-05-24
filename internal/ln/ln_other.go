//go:build !windows

package ln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// On non-Windows, ln creates hard links by default (POSIX behavior)
var otherCmd = &cli.Command{
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

func init() {
	Cmd = otherCmd
}
