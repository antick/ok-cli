# OK CLI v0.1.1 Release Notes

## Overview
OK CLI is a versatile command-line tool that provides various utilities for file operations, building Go programs, and managing Docker containers. This release introduces several powerful commands to enhance your development workflow.

## New Features

### Process Management
- **`ok kill [--port] <port>`** - Find and kill processes listening on a TCP port
  - Shows a table of processes (COMMAND, USER, PID, NAME) using the port
  - Prompts for confirmation before killing: `Proceed to kill them? [Y/n]:` (Enter defaults to Yes)
  - Uses `SIGKILL (-9)` to ensure processes exit
  - Accepts port as either a flag (`--port 3000`) or positional argument (`3000`)
  - If nothing is listening, prints a friendly message

### File Operations
- **`ok copy <source> [to] <destination>`** - Copy files or directories
  - Supports '~' in paths
  - You can omit the 'to' keyword
  - Example: `ok copy src.txt to ~/Downloads/`

- **`ok move <source> [to] <destination>`** - Move files or directories
  - Falls back to copy+delete across volumes
  - Supports '~' in paths
  - Example: `ok move ./bin to ~/bin`

- **`ok remove <file_or_directory> [--permanent|-p]`** - Remove files or directories
  - By default moves to Trash
  - Use `-p` or `--permanent` to permanently delete
  - Example: `ok remove ./dist --permanent`

### Development Tools
- **`ok build <input_file> [as/to] <output_file>`** - Build Go programs
  - Defaults output name to 'main' if not provided
  - Requires a .go file
  - Example: `ok build main.go as app`

- **`ok docker`** - Manage Docker containers
  - Launches an interactive UI to manage Docker containers

### Utility Commands
- **`ok version`** - Print the version number of OK CLI
- **`ok help`** - Show detailed help and usage information
  - Get an overview of all commands with examples

## Configuration
- Creates `~/.ok/config.yaml` file to set preferred defaults
- Supports verbose output with `-v` or `--verbose` flag

## System Requirements
- macOS (uses `lsof` for process management, available by default)
- Go 1.22+ (for building from source)

## Installation

### Option 1: Go Toolchain
```bash
go install github.com/antick/ok-cli@main

Ensure your `$GOPATH/bin (usually $HOME/go/bin)` is on `PATH`:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

### Option 2: Local Build
```bash
git clone https://github.com/antick/ok-cli.git
cd ok-cli
go build -o ok .
./ok help
```
