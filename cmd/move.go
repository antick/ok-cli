package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/antick/ok-cli/utils"
)

func HandleMove(cmd *cobra.Command, args []string) {
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

	err = utils.MoveFileOrDir(source, destination)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	if verbose {
		color.Green("Successfully moved %s to %s", source, destination)
	}
}
