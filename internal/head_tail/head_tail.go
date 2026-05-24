package head_tail

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

var Head = &cli.Command{
	Name:  "head",
	Usage: "Display the first lines of a file",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "lines", Aliases: []string{"n"}, Value: 10, Usage: "Number of lines"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		n := cmd.Int("lines")
		files := cmd.Args().Slice()

		if len(files) == 0 {
			return printLines(os.Stdin, n, true)
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("head: %s: %v", file, err)
			}
			if err := printLines(f, n, true); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
		return nil
	},
}

var Tail = &cli.Command{
	Name:  "tail",
	Usage: "Display the last lines of a file",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "lines", Aliases: []string{"n"}, Value: 10, Usage: "Number of lines"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		n := cmd.Int("lines")
		files := cmd.Args().Slice()

		if len(files) == 0 {
			return printLines(os.Stdin, n, false)
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("tail: %s: %v", file, err)
			}
			if err := printLines(f, n, false); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
		return nil
	},
}

func printLines(r io.Reader, n int, head bool) error {
	scanner := bufio.NewScanner(r)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if head {
			fmt.Println(line)
			n--
			if n <= 0 {
				return scanner.Err()
			}
		} else {
			lines = append(lines, line)
		}
	}

	if !head && len(lines) > 0 {
		start := 0
		if len(lines) > n {
			start = len(lines) - n
		}
		for _, line := range lines[start:] {
			fmt.Println(line)
		}
	}

	return scanner.Err()
}
