package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const Version = "0.1.1"

// HandleVersion prints the current version of the CLI.
func HandleVersion(cmd *cobra.Command, args []string) {
	color.Cyan("OK CLI Version: %s", Version)
	fmt.Println("A super CLI with super powers")
}
