package cp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "cp",
	Usage: "Copy files and directories",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "recursive", Aliases: []string{"r"}, Usage: "Copy directories recursively"},
		&cli.BoolFlag{Name: "interactive", Aliases: []string{"i"}, Usage: "Prompt before overwrite"},
		&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Print source/destination paths"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		recursive := cmd.Bool("recursive")
		interactive := cmd.Bool("interactive")
		verbose := cmd.Bool("verbose")

		files := cmd.Args().Slice()
		if len(files) < 2 {
			return fmt.Errorf("missing file operand")
		}

		source := files[:len(files)-1]
		dest := files[len(files)-1]

		destInfo, err := os.Stat(dest)
		if err != nil {
			if os.IsNotExist(err) {
				if len(source) != 1 {
					return fmt.Errorf("target '%s' is not a directory", dest)
				}
				dest = filepath.Join(filepath.Dir(dest), filepath.Base(source[0]))
			} else {
				return err
			}
		} else if destInfo.IsDir() {
			for _, s := range source {
				dest = filepath.Join(dest, filepath.Base(s))
			}
		}

		if len(source) == 1 {
			return copyFile(source[0], dest, recursive, interactive, verbose)
		}

		for _, s := range source {
			d := filepath.Join(dest, filepath.Base(s))
			if err := copyFile(s, d, recursive, interactive, verbose); err != nil {
				return err
			}
		}
		return nil
	},
}

func copyFile(src, dst string, recursive, interactive, verbose bool) error {
	if !canOverwrite(dst, interactive) {
		return nil
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if !recursive {
			return fmt.Errorf("omitting directory '%s'", src)
		}
		return copyDir(src, dst)
	}

	return copySingleFile(src, dst, verbose)
}

func copySingleFile(src, dst string, verbose bool) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, _ := os.Stat(src)
	if srcInfo != nil {
		os.Chmod(dst, srcInfo.Mode())
		os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
	}

	if verbose {
		fmt.Printf("%s -> %s\n", src, dst)
	}
	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copySingleFile(srcPath, dstPath, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func canOverwrite(dst string, interactive bool) bool {
	_, err := os.Stat(dst)
	if err != nil {
		return true
	}

	if !interactive {
		return true
	}

	fmt.Printf("overwrite '%s'? (y/N) ", dst)
	var answer string
	fmt.Scanln(&answer)
	return answer == "y" || answer == "Y"
}
