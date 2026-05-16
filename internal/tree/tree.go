package tree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "tree",
	Usage: "Display directory structure",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "depth, L", Value: 3, Usage: "Maximum depth"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		depth := cmd.Int("depth")
		dirs := cmd.Args().Slice()

		if len(dirs) == 0 {
			dirs = []string{"."}
		}

		for _, dir := range dirs {
			fmt.Println(dir)
			printTree(dir, "", depth, true)
		}
		return nil
	},
}

func printTree(dir, prefix string, depth int, isLast bool) {
	if depth <= 0 {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sortEntries(entries)

	for i, entry := range entries {
		isLastEntry := i == len(entries)-1
		connector := "├── "
		if isLastEntry {
			connector = "└── "
		}

		if entry.IsDir() {
			fmt.Printf("%s%s%s/\n", prefix, connector, entry.Name())
			newPrefix := prefix
			if isLastEntry {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(dir, entry.Name()), newPrefix, depth-1, isLastEntry)
		} else {
			fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
		}
	}
}

func sortEntries(entries []os.DirEntry) {
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].IsDir() && !entries[j].IsDir() {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}
