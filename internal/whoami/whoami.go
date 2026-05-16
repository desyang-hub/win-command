package whoami

import (
	"context"
	"fmt"
	"os/user"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "whoami",
	Usage: "Display current username",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		u, err := user.Current()
		if err != nil {
			return err
		}
		fmt.Println(u.Username)
		return nil
	},
}
