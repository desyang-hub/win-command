package zip

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v4"
	"github.com/urfave/cli/v3"
)

var ExtractCmd = &cli.Command{
	Name:  "unzip",
	Usage: "Extract a zip archive",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output directory (default: current)"},
		&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Show detailed output (deprecated: always on)"},
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

func extractZip(archive, dest string, _ bool) error {
	// 自动检测 ZIP 文件的编码：先试 UTF-8，再试 GBK
	encoding := detectZipEncoding(archive)

	z := archiver.Zip{
		TextEncoding: encoding,
	}

	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	archiveInfo, err := os.Stat(archive)
	if err != nil {
		return err
	}

	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	var extracted int
	fileCount := 0

	// 统计文件总数
	statF, _ := os.Open(archive)
	if statF != nil {
		defer statF.Close()
		statInfo, _ := statF.Stat()
		if statInfo != nil {
			zr, _ := zip.NewReader(statF, statInfo.Size())
			if zr != nil {
				fileCount = len(zr.File)
			}
		}
	}

	fmt.Printf("Extracting: %s (%.2f KB, %d entries)\n", archiveInfo.Name(), float64(archiveInfo.Size())/1024, fileCount)
	fmt.Printf("Destination: %s\n\n", destAbs)

	err = z.Extract(context.Background(), f, func(ctx context.Context, file archiver.FileInfo) error {
		cleanName := filepath.Clean(filepath.Join(destAbs, file.NameInArchive))
		cleanDest := filepath.Clean(destAbs) + string(filepath.Separator)
		if !strings.HasPrefix(cleanName, cleanDest) || cleanName == filepath.Clean(destAbs) {
			return fmt.Errorf("illegal path: %s", file.NameInArchive)
		}

		if file.IsDir() {
			if err := os.MkdirAll(cleanName, file.Mode()); err != nil {
				return err
			}
			fmt.Printf("  [DIR]  %s\n", cleanName)
			extracted++
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(cleanName), 0755); err != nil {
			return err
		}

		opened, err := file.Open()
		if err != nil {
			return err
		}
		defer opened.Close()

		outFile, err := os.OpenFile(cleanName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		written, err := io.Copy(outFile, opened)
		if err != nil {
			return err
		}

		size := fmt.Sprintf("%d B", written)
		if written > 1024*1024 {
			size = fmt.Sprintf("%.2f MB", float64(written)/1024/1024)
		} else if written > 1024 {
			size = fmt.Sprintf("%.1f KB", float64(written)/1024)
		}

		fmt.Printf("  [FILE] %s (%s)\n", cleanName, size)
		extracted++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nDone! Extracted %d files to: %s\n", extracted, destAbs)
	return nil
}

// detectZipEncoding scans a ZIP file and picks the best text encoding.
// Uses klauspost/compress/zip which correctly handles the UTF-8 flag (bit 11).
// Falls back to GBK for Chinese Windows ZIP files without the UTF-8 flag.
func detectZipEncoding(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "gbk"
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "gbk"
	}

	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return "gbk"
	}

	for _, file := range zr.File {
		if !utf8.ValidString(file.Name) {
			return "gbk"
		}
		if file.FileHeader.NonUTF8 && hasNonASCII(file.Name) {
			return "gbk"
		}
	}

	return ""
}

func hasNonASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return true
		}
	}
	return false
}
