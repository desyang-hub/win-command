package zip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

var ExtractCmd = &cli.Command{
	Name:  "unzip",
	Usage: "Extract a zip archive",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output, o", Usage: "Output directory (default: current)"},
		&cli.BoolFlag{Name: "verbose, v", Usage: "Show extracted files"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		output := cmd.String("output")
		verbose := cmd.Bool("verbose")

		files := cmd.Args().Slice()
		if len(files) < 1 {
			return fmt.Errorf("usage: win unzip [-o output] archive.zip")
		}

		archive := files[0]
		if output == "" {
			output = "."
		}

		return extractZip(archive, output, verbose)
	},
}

func extractZip(archive, dest string, verbose bool) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := extractFile(reader, file, dest, verbose); err != nil {
			return err
		}
	}

	fmt.Printf("Extracted %d files to: %s\n", len(reader.File), dest)
	return nil
}

func extractFile(reader *zip.ReadCloser, file *zip.File, dest string, verbose bool) error {
	cleanName := filepath.Join(dest, file.Name)
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	if !strings.HasPrefix(filepath.Clean(cleanName), cleanDest) && cleanName != cleanDest {
		return fmt.Errorf("illegal path: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		os.MkdirAll(cleanName, file.Mode())
		if verbose {
			fmt.Printf("  %s/\n", cleanName)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cleanName), 0755); err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(cleanName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return err
	}

	outFile.Chmod(file.Mode())

	if verbose {
		fmt.Printf("  %s\n", file.Name)
	}
	return nil
}
