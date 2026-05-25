//go:build windows

package ln

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/urfave/cli/v3"
)

func createShortcut(target, linkName string) error {
	// Ensure linkName ends with .lnk
	if !strings.HasSuffix(strings.ToLower(linkName), ".lnk") {
		linkName = linkName + ".lnk"
	}

	psScript := `$shell = New-Object -ComObject WScript.Shell
$sc = $shell.CreateShortcut($args[0])
$sc.TargetPath = $args[1]
$sc.WorkingDirectory = $args[2]
$sc.Save()`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript, linkName, target, filepath.Dir(target))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create shortcut '%s' -> '%s': %s", linkName, target, strings.TrimSpace(string(output)))
	}
	return nil
}

// selfElevate restarts the current process with administrator privileges via PowerShell.
// It is called when symlink creation fails due to insufficient permissions.
func selfElevate() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Note: elevation failed (%v), try running as administrator manually.\n", err)
		os.Exit(1)
	}

	argsStr := marshalPSArray(os.Args[1:])

	psScript := "Start-Process '" + exe + "' -ArgumentList @(" + argsStr + ") -Verb RunAs"

	cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Note: elevation failed (%v), try running as administrator manually.\n", err)
	}
	os.Exit(0)
}

// marshalPSArray builds a PowerShell array literal from Go args ('arg1','arg2',...).
func marshalPSArray(args []string) string {
	var parts []string
	for _, arg := range args {
		escaped := strings.ReplaceAll(arg, "'", "''")
		parts = append(parts, "'"+escaped+"'")
	}
	return strings.Join(parts, ",")
}

// isSymlinkPermissionError checks if the error is a Windows symlink permission denial.
func isSymlinkPermissionError(err error) bool {
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "access is denied") {
		return true
	}
	if strings.Contains(msg, "required privilege is not held") {
		return true
	}
	if strings.Contains(msg, "0x5") {
		return true
	}
	// Also check for syscall error with ACCESS_DENIED code
	var se *os.SyscallError
	if errors.As(err, &se) && se.Err == syscall.EACCES {
		return true
	}
	return false
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
			if err := createSymlink(target, linkName); err != nil {
				if isSymlinkPermissionError(err) {
					fmt.Fprintf(os.Stderr, "%v\n  Attempting UAC elevation...\n", err)
					selfElevate()
				}
				return err
			}
			return nil
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
