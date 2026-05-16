package main

import (
	"context"
	"fmt"
	"os"

	"github.com/desyang-hub/win-command/internal/cat"
	"github.com/desyang-hub/win-command/internal/clear"
	"github.com/desyang-hub/win-command/internal/cp"
	"github.com/desyang-hub/win-command/internal/date"
	"github.com/desyang-hub/win-command/internal/echo"
	"github.com/desyang-hub/win-command/internal/head_tail"
	"github.com/desyang-hub/win-command/internal/ln"
	"github.com/desyang-hub/win-command/internal/ls"
	"github.com/desyang-hub/win-command/internal/mkdir"
	"github.com/desyang-hub/win-command/internal/mv"
	"github.com/desyang-hub/win-command/internal/pwd"
	"github.com/desyang-hub/win-command/internal/rm"
	"github.com/desyang-hub/win-command/internal/tree"
	"github.com/desyang-hub/win-command/internal/which"
	"github.com/desyang-hub/win-command/internal/whoami"
	"github.com/desyang-hub/win-command/internal/zip"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "win",
		Usage: "Linux-style commands for Windows",
		Commands: []*cli.Command{
			// Core commands
			rm.Cmd,
			cp.Cmd,
			mv.Cmd,
			ln.Cmd,
			zip.Cmd,
			zip.ExtractCmd,

			// File viewing
			cat.Cmd,
			head_tail.Head,
			head_tail.Tail,
			tree.Cmd,

			// System utilities
			clear.Cmd,
			which.Cmd,
			date.Cmd,
			whoami.Cmd,

			// Directory
			mkdir.Cmd,
			pwd.Cmd,
			ls.Cmd,

			// Text
			echo.Cmd,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
