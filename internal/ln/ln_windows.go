//go:build windows

package ln

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func createShortcut(target, linkName string) error {
	// Ensure linkName ends with .lnk
	if !strings.HasSuffix(strings.ToLower(linkName), ".lnk") {
		linkName = linkName + ".lnk"
	}

	// Escape single quotes for PowerShell
	et := strings.ReplaceAll(target, "'", "''")
	el := strings.ReplaceAll(linkName, "'", "''")
	edir := strings.ReplaceAll(filepath.Dir(target), "'", "''")

	psScript := fmt.Sprintf(
		"$shell = New-Object -ComObject WScript.Shell; "+
			"$sc = $shell.CreateShortcut('%s'); "+
			"$sc.TargetPath = '%s'; "+
			"$sc.WorkingDirectory = '%s'; "+
			"$sc.Save()",
		el, et, edir,
	)

	cmd := exec.Command("powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create shortcut '%s' -> '%s': %s", linkName, target, strings.TrimSpace(string(output)))
	}
	return nil
}

// Windows-specific command definition
var winCmd = &cli.Command{
	Name:  "ln",
	Usage: "Make links between files (shortcut on Windows, hard link on other platforms)",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "symbolic", Aliases: []string{"s"}, Usage: "Create symbolic link (requires admin/Developer Mode)"},
		&cli.BoolFlag{Name: "hard", Usage: "Create hard link instead of shortcut"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		symbolic := cmd.Bool("symbolic")
		hardLink := cmd.Bool("hard")

		files := cmd.Args().Slice()
		if len(files) < 2 {
			return fmt.Errorf("missing file operand")
		}

		target, _ := filepath.Abs(files[0])
		linkName, _ := filepath.Abs(files[1])

		if symbolic {
			return createSymlink(target, linkName)
		}
		if hardLink {
			return createHardlink(target, linkName)
		}
		// Default on Windows: shortcut (.lnk) — check target exists first
		if _, err := os.Stat(target); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("target '%s' does not exist", target)
			}
		}
		return createShortcut(target, linkName)
	},
}

func init() {
	Cmd = winCmd
}
