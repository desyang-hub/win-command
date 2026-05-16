package echo

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "echo",
	Usage: "Display a line of text",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "n", Usage: "Do not output trailing newline"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		noNewline := cmd.Bool("n")
		fmt.Fprint(os.Stdout, strings.Join(cmd.Args().Slice(), " "))
		if !noNewline {
			fmt.Fprintln(os.Stdout)
		}
		return nil
	},
}
