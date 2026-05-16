package clear

import (
	"context"
	"fmt"
	"runtime"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "clear",
	Usage: "Clear the terminal screen",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if runtime.GOOS == "windows" {
			fmt.Print("\x1b[2J\x1b[H")
		} else {
			fmt.Print("\033[2J\033[H")
		}
		return nil
	},
}
