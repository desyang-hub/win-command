package mv

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "mv",
	Usage: "Move/rename files or directories",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "interactive, i", Usage: "Prompt before overwrite"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		interactive := cmd.Bool("interactive")

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

		for _, s := range source {
			if err := moveFile(s, dest, interactive); err != nil {
				return err
			}
		}
		return nil
	},
}

func moveFile(src, dst string, interactive bool) error {
	if !canOverwrite(dst, interactive) {
		return nil
	}

	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)

	if sameDrive(src, dst) {
		if err := os.Rename(src, dst); err == nil {
			return nil
		}
	}

	if err := copyFile(src, dst); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		os.Remove(dst)
		return err
	}
	return nil
}

func sameDrive(src, dst string) bool {
	if runtime.GOOS != "windows" {
		return true
	}
	return filepath.VolumeName(src) == filepath.VolumeName(dst)
}

func copyFile(src, dst string) error {
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
