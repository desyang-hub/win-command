package date

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "date",
	Usage: "Display current date and time",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Value: "", Usage: "Custom format string"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		format := cmd.String("format")
		if format == "" {
			format = "2006-01-02 15:04:05 Mon Jan"
		}
		fmt.Println(time.Now().Format(format))
		return nil
	},
}
