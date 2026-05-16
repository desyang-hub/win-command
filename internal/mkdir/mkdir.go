package mkdir

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "mkdir",
	Usage: "Create directories",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "parents, p", Usage: "Create parent directories as needed"},
		&cli.BoolFlag{Name: "verbose, v", Usage: "Print each created directory"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		parents := cmd.Bool("parents")
		verbose := cmd.Bool("verbose")

		dirs := cmd.Args().Slice()
		if len(dirs) == 0 {
			return fmt.Errorf("missing directory operand")
		}

		for _, dir := range dirs {
			var err error
			if parents {
				err = os.MkdirAll(dir, 0755)
			} else {
				err = os.Mkdir(dir, 0755)
			}

			if err != nil {
				return fmt.Errorf("mkdir: %s: %v", dir, err)
			}

			if verbose {
				fmt.Printf("created directory '%s'\n", dir)
			}
		}
		return nil
	},
}
