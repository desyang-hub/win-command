package rm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:  "rm",
	Usage: "Remove files or directories (moves to recycle bin by default)",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "recursive", Aliases: []string{"r"}, Usage: "Remove directories recursively"},
		&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "Ignore nonexistent files, never prompt"},
		&cli.BoolFlag{Name: "permanent", Aliases: []string{"P"}, Usage: "Permanently delete (skip recycle bin)"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		recursive := cmd.Bool("recursive")
		force := cmd.Bool("force")
		permanent := cmd.Bool("permanent")

		files := cmd.Args().Slice()
		if len(files) == 0 {
			return fmt.Errorf("missing file operand")
		}

		for _, file := range files {
			if err := remove(file, recursive, force, permanent); err != nil {
				return err
			}
		}
		return nil
	},
}

func remove(file string, recursive, force, permanent bool) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && force {
			return nil
		}
		return err
	}

	if info.IsDir() && !recursive {
		if !force {
			return fmt.Errorf("cannot remove '%s': it is a directory (use -r)", path)
		}
		return nil
	}

	if permanent {
		return deletePermanent(path, info.IsDir())
	}
	return moveToTrash(path, info.IsDir())
}

func moveToTrash(path string, isDir bool) error {
	if isDir {
		return moveToTrashRecursively(path)
	}
	return moveToRecycleBin(path)
}

func moveToRecycleBin(path string) error {
	dll := syscall.NewLazyDLL("shell32.dll")
	proc := dll.NewProc("SHFileOperationW")

	pFrom, err := makeDoubleNullStr(path)
	if err != nil {
		return fmt.Errorf("path encoding failed: %w", err)
	}

	var op shOp
	op.wFunc = FO_DELETE
	op.pFrom = uintptr(unsafe.Pointer(&pFrom[0]))
	op.fFlags = FOF_ALLOWUNDO | FOF_NOERRORSDIALOG
	op.pTo = 0

	ret, _, lastErr := proc.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		errStr := ""
		if lastErr != nil {
			errStr = lastErr.Error()
		}
		return fmt.Errorf("failed to recycle '%s': SHFileOperationW returned %d%s", path, ret, errStr)
	}
	return nil
}

func moveToTrashRecursively(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		isDir := entry.IsDir()
		if err := moveToTrash(entryPath, isDir); err != nil {
			return err
		}
	}
	return moveToRecycleBin(path)
}

func makeDoubleNullStr(s string) ([]uint16, error) {
	r, err := syscall.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	return append(r, 0), nil
}

func deletePermanent(path string, isDir bool) error {
	if isDir {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

// Windows Shell API constants
const (
	FO_DELETE          = 0x0003
	FOF_ALLOWUNDO      = 0x0040 // CRITICAL: must be set to recycle the file!
	FOF_NOCONFIRMATION = 0x0010
	FOF_NOERRORSDIALOG = 0x1000
)

// shOp matches SHFILEOPSTRUCTW memory layout.
// NOTE: 64-bit only — pointer sizes differ on 386. Do not build for windows/386.
type shOp struct {
	hwnd      uintptr
	wFunc     uint32
	pFrom     uintptr
	pTo       uintptr
	fFlags    uint16
	fAborted  byte
	_         [3]byte
	fileName  uintptr
	fAborted2 byte
	_         [7]byte
}
