package cmd

import (
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/antick/ok-cli/utils"
)

func HandleBuild(cmd *cobra.Command, args []string) {
	inputFile, outputFile, err := utils.ParseSourceAndDestination(args)
	if err != nil {
		color.Red("Error: %v", err)
		cmd.Usage()
		return
	}

	if outputFile == "" {
		outputFile, _ = cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = "main"
		}
	}

	if filepath.Ext(inputFile) != ".go" {
		color.Red("Error: Unsupported file type: %s", inputFile)
		return
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	buildGoProgram(inputFile, outputFile, verbose)
}

func buildGoProgram(inputFile, outputFile string, verbose bool) {
	cmd := exec.Command("go", "build", "-o", outputFile, inputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Error building Go program: %v", err)
		color.Yellow(string(output))
		return
	}

	if verbose {
		color.Green("Successfully built %s as %s", inputFile, outputFile)
	}
}
