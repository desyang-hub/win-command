package clear

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "clear",
	Usage: "Clear the terminal screen",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Print("\033[2J\033[H")
		return nil
	},
}
