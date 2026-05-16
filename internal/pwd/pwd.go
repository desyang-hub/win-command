package pwd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "pwd",
	Usage: "Print name of current working directory",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(filepath.ToSlash(dir))
		return nil
	},
}
