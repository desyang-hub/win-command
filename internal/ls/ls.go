package ls

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "ls",
	Usage: "List directory contents",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "long, l", Usage: "Long listing format"},
		&cli.BoolFlag{Name: "all, a", Usage: "Include hidden files"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		long := cmd.Bool("long")
		all := cmd.Bool("all")

		files := cmd.Args().Slice()
		if len(files) == 0 {
			files = []string{"."}
		}

		for _, file := range files {
			if long {
				listLong(file, all)
			} else {
				listShort(file, all)
			}
		}
		return nil
	},
}

func listShort(path string, all bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		fmt.Println(entry.Name())
	}
	return nil
}

func listLong(path string, all bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, entry := range entries {
		if !all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, _ := entry.Info()
		size := int64(0)
		if info != nil {
			size = info.Size()
		}
		mode := "file"
		if entry.Type().IsDir() {
			mode = "dir"
		}
		fmt.Fprintf(w, "%s  %10d  %s\n", mode, size, entry.Name())
	}
	w.Flush()
	return nil
}
