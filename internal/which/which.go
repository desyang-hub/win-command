package which

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "which",
	Usage: "Show the full path of a command",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		files := cmd.Args().Slice()
		if len(files) == 0 {
			return fmt.Errorf("missing command name")
		}

		for _, name := range files {
			path, err := findCommand(name)
			if err != nil {
				return fmt.Errorf("which: no %s in PATH", name)
			}
			fmt.Println(path)
		}
		return nil
	},
}

func findCommand(name string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name += ".exe"
	}

	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		fullPath := filepath.Join(dir, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("not found")
}
