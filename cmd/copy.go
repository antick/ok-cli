package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/antick/ok-cli/utils"
)

func HandleCopy(cmd *cobra.Command, args []string) {
	source, destination, err := utils.ParseSourceAndDestination(args)
	if err != nil {
		color.Red("Error: %v", err)
		cmd.Usage()
		return
	}

	if destination == "" {
		destination, _ = cmd.Flags().GetString("destination")
		if destination == "" {
			color.Red("Error: No destination specified")
			cmd.Usage()
			return
		}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	err = utils.CopyFileOrDir(source, destination)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	if verbose {
		color.Green("Successfully copied %s to %s", source, destination)
	}
}
