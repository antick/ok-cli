package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/antick/ok/utils"
)

func HandleRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		color.Red("Error: No file or directory specified")
		cmd.Usage()
		return
	}

	permanent, _ := cmd.Flags().GetBool("permanent")
	verbose, _ := cmd.Flags().GetBool("verbose")

	for _, path := range args {
		path = utils.ExpandPath(path)
		err := removeFileOrDir(path, permanent)
		if err != nil {
			color.Red("Error removing %s: %v", path, err)
		} else if verbose {
			if permanent {
				color.Green("Successfully deleted %s", path)
			} else {
				color.Green("Successfully moved %s to trash", path)
			}
		}
	}
}

func removeFileOrDir(path string, permanent bool) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}

	if permanent {
		return os.RemoveAll(path)
	}

	// Move to system trash
	err = utils.MoveToTrash(path)
	if err != nil {
		return fmt.Errorf("error moving to trash: %w", err)
	}

	return nil
}
