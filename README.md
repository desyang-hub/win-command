# win-command

Linux-style CLI tools for Windows, written in Go.

> [English](README.md) | [中文](README_zh.md)

## Overview

Bring familiar Linux commands to Windows. Solves the gap between Windows built-in commands and POSIX tools in usability and consistency. All commands are prefixed with `win` and powered by native Windows Shell APIs.

## Why

| Scenario | Windows Default | win-command |
|----------|----------------|-------------|
| Delete files | Permanently deleted | Moved to Recycle Bin (recoverable) |
| Directory tree | `dir` shows flat list | `tree` shows visual hierarchy |
| Clear screen | `cls` (non-POSIX) | `clear` (standard ANSI escape) |
| Find executables | Requires PowerShell | `which` finds PATH binaries |
| Compress files | No native CLI | `zip`/`unzip` from terminal |

## Installation

### Option 1: Download prebuilt binary

Download the latest `win.exe` from [Releases](https://github.com/desyang-hub/win-command/releases) and place it in a directory on your `PATH` (e.g. `C:\Windows`).

### Option 2: Build from source

```bash
git clone https://github.com/desyang-hub/win-command.git
cd win-command
go build -o win.exe ./win/main.go
```

### Option 3: Go install

```bash
go install github.com/desyang-hub/win-command/win@latest
```

## Command Reference

### Core Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `rm` | Remove files/directories (moves to Recycle Bin by default) | `-r` recursive · `-f` force · `-P` permanent delete |
| `cp` | Copy files/directories | `-r` recursive · `-i` interactive · `-v` verbose |
| `mv` | Move/rename | `-i` interactive |
| `ln` | Create links | `-s` symbolic link |
| `zip` | Create zip archive | `-r` recursive |
| `unzip` | Extract zip archive | `-o` output dir · `-v` verbose |

### File Viewing

| Command | Description | Flags |
|---------|-------------|-------|
| `cat` | Display file contents | Supports stdin pipe |
| `head` | Show file header | `-n` lines (default 10) |
| `tail` | Show file tail | `-n` lines (default 10) |
| `tree` | Show directory tree | `-L` max depth (default 3) |

### System Utilities

| Command | Description |
|---------|-------------|
| `clear` | Clear terminal screen |
| `which` | Show full path of a command |
| `date` | Display current date and time |
| `whoami` | Show current username |
| `mkdir` | Create directories |
| `pwd` | Print working directory |
| `ls` | List directory contents |
| `echo` | Print text |
| `echo -n` | Print without trailing newline |

## Usage Examples

```bash
# Delete file (moves to Recycle Bin, recoverable)
win rm file.txt

# Delete directory recursively
win rm -r my_folder

# Permanent delete (skips Recycle Bin)
win rm -P sensitive_file.txt

# Copy directory
win cp -r source_dir dist_dir

# Move file (prompt before overwrite)
win mv -i old_name.txt new_name.txt

# Create symbolic link
win ln -s target.txt link.txt

# Compress files
win zip -r archive.zip folder/

# Extract (with verbose output)
win unzip -v archive.zip

# Show directory tree
win tree -L 2 src/

# Show last 20 lines of a file
win tail -n 20 logfile.txt

# Clear screen
win clear

# Find where git is installed
win which git
```

## Architecture

```
win-command/
├── cmd/main.go                  # Entry point, CLI registration
├── internal/
│   ├── rm/rm.go                 # Recycle Bin delete (SHFileOperationW)
│   ├── cp/cp.go                 # Recursive copy
│   ├── mv/mv.go                 # Move/rename
│   ├── ln/ln.go                 # Symlink/hard link
│   ├── zip/                     # zip/unzip archiving
│   ├── cat/cat.go               # File content viewer
│   ├── head_tail/head_tail.go   # head/tail
│   ├── tree/tree.go             # Directory tree
│   ├── clear/clear.go           # Screen clearing
│   ├── which/which.go           # Command lookup
│   ├── date/date.go             # Date/time
│   ├── whoami/whoami.go         # Username
│   ├── mkdir/mkdir.go           # Directory creation
│   ├── pwd/pwd.go               # Working directory
│   ├── ls/ls.go                 # Directory listing
│   └── echo/echo.go             # Text output
└── go.mod
```

### Technology Stack

| Component | Choice | Rationale |
|-----------|--------|-----------|
| CLI framework | urfave/cli/v3 | Lightweight, minimal boilerplate |
| Recycle Bin | Shell32.SHFileOperationW | Stable across all Windows versions |
| Compression | Go stdlib `archive/zip` | Zero external dependencies |

### Recycle Bin Implementation

The `rm` command uses the `SHFileOperationW` API:

- **FOF_ALLOWUNDO** — preserves undo capability (key flag for recycling)
- **FOF_NOERRORSDIALOG** — suppresses error dialogs
- **FOF_NOCONFIRMATION** — no confirmation prompt

Directories are processed recursively, each file gets its own undo record in the Recycle Bin.

## Build & Release

```bash
# Local build
go build -o win.exe ./win/main.go

# Cross-compile Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o win-amd64.exe ./win/main.go

# Cross-compile Windows ARM64
GOOS=windows GOARCH=arm64 go build -o win-arm64.exe ./win/main.go
```

CI/CD automatically builds platform binaries on push and PR, and publishes tagged releases.

## License

[MIT](LICENSE)
