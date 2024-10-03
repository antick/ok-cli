package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// MoveToTrash moves a file or directory to the system trash
func MoveToTrash(path string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("moving to trash is currently only supported on macOS")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	script := fmt.Sprintf(`
		on run {p}
			tell application "Finder"
				set p to POSIX file p as alias
				delete p
			end tell
		end run
	`)

	cmd := exec.Command("osascript", "-e", script, absPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error moving to trash: %s", output)
	}

	return nil
}
