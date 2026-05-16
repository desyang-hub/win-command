# win-command

Windows 版的 Linux 风格命令行工具，用 Go 实现。

> [English](README.md) | 中文

## 概述

在 Windows 上提供类似 Linux 的常用命令，解决 Windows 原生命令在易用性、功能一致性上的不足。所有命令统一以 `win` 前缀调用，底层使用 Windows Shell API 实现。

## 为什么需要

| 场景 | Windows 默认 | win-command |
|------|-------------|-------------|
| 删除文件 | 直接永久删除 | 默认移入回收站，可恢复 |
| 目录结构 | `dir` 只显示列表 | `tree` 可视化树状结构 |
| 清屏 | `cls` 非 POSIX | `clear` 标准 ANSI 转义 |
| 文件搜索 | 需要 PowerShell 脚本 | `which` 快速查找可执行文件 |
| 压缩包 | 无原生 CLI | `zip`/`unzip` 命令行操作 |

## 安装

### 方式一：下载预编译版本

从 [Releases](https://github.com/desyang-hub/win-command/releases) 下载最新 `win.exe`，放到 PATH 中的某个目录（如 `C:\Windows`）。

### 方式二：从源码构建

```bash
git clone https://github.com/desyang-hub/win-command.git
cd win-command
go build -o win.exe ./cmd/main.go
```

### 方式三：Go install

```bash
go install github.com/desyang-hub/win-command/cmd@latest
```

## 命令参考

### 核心命令

| 命令 | 说明 | 选项 |
|------|------|------|
| `rm` | 删除文件/目录（默认移入回收站） | `-r` 递归 · `-f` 强制 · `-P` 永久删除 |
| `cp` | 复制文件/目录 | `-r` 递归 · `-i` 交互 · `-v` 详细 |
| `mv` | 移动/重命名 | `-i` 交互 |
| `ln` | 创建链接 | `-s` 符号链接 |
| `zip` | 创建压缩包 | `-r` 递归 |
| `unzip` | 解压压缩包 | `-o` 输出目录 · `-v` 详细 |

### 文件查看

| 命令 | 说明 | 选项 |
|------|------|------|
| `cat` | 查看文件内容 | 支持管道输入 |
| `head` | 查看文件头部 | `-n` 行数（默认10） |
| `tail` | 查看文件尾部 | `-n` 行数（默认10） |
| `tree` | 显示目录树 | `-L` 最大深度（默认3） |

### 系统工具

| 命令 | 说明 |
|------|------|
| `clear` | 清屏 |
| `which` | 查找可执行文件路径 |
| `date` | 显示日期时间 |
| `whoami` | 显示当前用户名 |
| `mkdir` | 创建目录 |
| `pwd` | 显示当前工作目录 |
| `ls` | 列出目录内容 |
| `echo` | 输出文本 |
| `echo -n` | 不输出换行 |

## 使用示例

```bash
# 删除文件（移入回收站，可恢复）
win rm file.txt

# 递归删除目录
win rm -r my_folder

# 永久删除（不走回收站）
win rm -P sensitive_file.txt

# 复制目录
win cp -r source_dir dist_dir

# 移动文件（交互，覆盖时提示）
win mv -i old_name.txt new_name.txt

# 创建符号链接
win ln -s target.txt link.txt

# 压缩文件
win zip -r archive.zip folder/

# 解压（显示进度）
win unzip -v archive.zip

# 查看目录树
win tree -L 2 src/

# 查看文件末尾 20 行
win tail -n 20 logfile.txt

# 清屏
win clear

# 查找 git 的位置
win which git
```

## 技术细节

### 架构

```
win-command/
├── cmd/main.go                  # 入口，CLI 注册
├── internal/
│   ├── rm/rm.go                 # 回收站删除（SHFileOperationW）
│   ├── cp/cp.go                 # 递归复制
│   ├── mv/mv.go                 # 移动/重命名
│   ├── ln/ln.go                 # 符号链接/硬链接
│   ├── zip/                     # zip/unzip 打包
│   ├── cat/cat.go               # 文件内容查看
│   ├── head_tail/head_tail.go   # head/tail
│   ├── tree/tree.go             # 目录树
│   ├── clear/clear.go           # 清屏
│   ├── which/which.go           # 查找命令
│   ├── date/date.go             # 日期
│   ├── whoami/whoami.go         # 用户名
│   ├── mkdir/mkdir.go           # 创建目录
│   ├── pwd/pwd.go               # 工作目录
│   ├── ls/ls.go                 # 列表
│   └── echo/echo.go             # 输出文本
└── go.mod
```

### 关键技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| CLI 框架 | urfave/cli/v3 | 轻量、代码少 |
| 回收站 | Shell32.SHFileOperationW | 所有 Windows 版本稳定支持 |
| 压缩 | Go 标准库 archive/zip | 零额外依赖 |

### 回收站实现

`rm` 默认使用 `SHFileOperationW` API 实现回收站功能：

- **FOF_ALLOWUNDO** — 保留撤销能力（移入回收站的关键 flag）
- **FOF_NOERRORSDIALOG** — 静默处理，不弹错误对话框
- **FOF_NOCONFIRMATION** — 不弹出确认对话框

目录采用递归逐个移入回收站，保证每个文件都有独立的撤销记录。

## 构建与发布

```bash
# 本地构建
go build -o win.exe ./cmd/main.go

# 交叉编译 Windows 64位
GOOS=windows GOARCH=amd64 go build -o win-amd64.exe ./cmd/main.go

# 交叉编译 Windows ARM64
GOOS=windows GOARCH=arm64 go build -o win-arm64.exe ./cmd/main.go
```

CI/CD 自动在推送和 PR 时构建所有平台的二进制文件并上传到 Release。

## 许可证

MIT
