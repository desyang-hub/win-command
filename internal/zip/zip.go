package zip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "zip",
	Usage: "Create a zip archive",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "recursive", Aliases: []string{"r"}, Usage: "Include directories recursively"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		recursive := cmd.Bool("recursive")

		files := cmd.Args().Slice()
		if len(files) < 2 {
			return fmt.Errorf("usage: win zip [-r] output.zip file1 [file2 ...]")
		}

		output := files[0]
		sources := files[1:]

		return createZip(output, sources, recursive)
	},
}

func createZip(output string, sources []string, recursive bool) error {
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := zip.NewWriter(outFile)
	defer writer.Close()

	for _, source := range sources {
		if err := addToFileWriter(writer, source, "", recursive); err != nil {
			return err
		}
	}

	fmt.Printf("Created: %s\n", output)
	return nil
}

func addToFileWriter(w *zip.Writer, source, base string, recursive bool) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("skipping directory, use -r: %s", source)
		}
		return addDirToZip(w, source, base)
	}

	return addFileToZip(w, source, base)
}

func addDirToZip(w *zip.Writer, dirPath, base string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dirPath, entry.Name())
		newBase := filepath.Join(base, entry.Name())

		if entry.IsDir() {
			if err := addDirToZip(w, path, newBase); err != nil {
				return err
			}
		} else {
			if err := addFileToZip(w, path, newBase); err != nil {
				return err
			}
		}
	}
	return nil
}

func addFileToZip(w *zip.Writer, path, base string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = base
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
