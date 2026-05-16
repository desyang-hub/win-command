package cat

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "cat",
	Usage: "Display file contents",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		files := cmd.Args().Slice()
		if len(files) == 0 {
			_, err := io.Copy(os.Stdout, os.Stdin)
			return err
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("cat: %s: %v", file, err)
			}
			if _, err := io.Copy(os.Stdout, f); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
		return nil
	},
}
