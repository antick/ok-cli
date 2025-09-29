# ok - CLI (v0.1.1 - alpha)
ok - a super CLI with super powers

## Installation

macOS options:

- Go toolchain (requires Go 1.22+):
  ```bash
  go install github.com/antick/ok-cli@latest
  ```
  Ensure your `$GOPATH/bin` (usually `$HOME/go/bin`) is on `PATH`:
  ```bash
  echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
  ```

- Local build:
  ```bash
  git clone https://github.com/antick/ok-cli.git
  cd ok-cli
  go build -o ok .
  ./ok help
  ```

- From GitHub Releases (when available):
  1. Download the latest `ok` binary for macOS (arm64 or amd64) from the Releases page.
  2. Make it executable and place it on your PATH:
     ```bash
     chmod +x ok && mv ok /usr/local/bin/ok
     ```

## Usage

Get an overview and examples:
```bash
ok help
```

Common commands:
```bash
ok copy <source> [to] <destination>
ok build <input_file> [as/to] <output_file>
ok move <source> [to] <destination>
ok remove <file_or_directory> [--permanent|-p]
ok docker
ok kill [--port] <port>
```

### Kill processes on a port (macOS)

Find and kill processes listening on a TCP port. This uses `lsof` under the hood (equivalent to `lsof -iTCP:<port> -sTCP:LISTEN`), shows the list of matching processes, and asks for confirmation before killing with `SIGKILL (-9)`.

Examples:
```bash
# Find what is using port 3000 and kill after confirmation
ok kill --port 3000

# You can also specify the port as a positional argument
ok kill 3000

# If you omit the port, the help menu is shown
ok kill
```

Behavior:
- Shows a table of processes (COMMAND, USER, PID, NAME) using the port.
- Prompts: `Proceed to kill them? [Y/n]:` Enter defaults to Yes.
- Uses `SIGKILL` to ensure processes exit.
- If nothing is listening, it prints a friendly message.

Notes:
- Requires `lsof` (available by default on macOS).
- You may need elevated privileges to kill some processes.

## Config

It creates ~/.ok/config.yaml file to set preferred defaults.

## Update

- If installed via Go:
  ```bash
  go install github.com/antick/ok-cli@latest
  ```
  This fetches and installs the latest version.
- If installed from GitHub Releases: re-download the latest binary and replace the existing `/usr/local/bin/ok`.

## License

GNU General Public License v3.0
