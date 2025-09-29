package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// HandleHelp prints a curated help for the root command.
func HandleHelp(c *cobra.Command, args []string) {
	root := c.Root()

	color.Cyan("OK CLI - a super CLI with super powers")
	fmt.Println()

	color.Yellow("Usage:")
	fmt.Println("  ok [command] [flags] [...args]")
	fmt.Println("  ok help")
	fmt.Println()

	color.Yellow("Available Commands:")
	// Collect and sort by name for stable output
	cmds := root.Commands()
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name() < cmds[j].Name() })
	for _, sc := range cmds {
		if sc.Hidden {
			continue
		}
		// Print primary usage and short description
		fmt.Printf("  %-24s %s\n", sc.UseLine(), sc.Short)
	}
	fmt.Println()

	color.Yellow("Global Flags:")
	fmt.Println("  -v, --verbose           verbose output")
	fmt.Println()

	color.Yellow("Detailed Usage:")
	fmt.Println("  ok copy <source> [to] <destination>")
	fmt.Println("    Copies files or directories. Supports '~' in paths. You can omit the 'to' keyword.")
	fmt.Println("    Example: ok copy src.txt to ~/Downloads/")
	fmt.Println()
	fmt.Println("  ok move <source> [to] <destination>")
	fmt.Println("    Moves files or directories. Falls back to copy+delete across volumes. Supports '~'.")
	fmt.Println("    Example: ok move ./bin to ~/bin")
	fmt.Println()
	fmt.Println("  ok build <input_file> [as/to] <output_file>")
	fmt.Println("    Builds Go programs. Defaults output name to 'main' if not provided. Requires a .go file.")
	fmt.Println("    Example: ok build main.go as app")
	fmt.Println()
	fmt.Println("  ok remove <file_or_directory> [-p|--permanent]")
	fmt.Println("    Removes files or directories. By default moves to Trash; use -p to permanently delete.")
	fmt.Println("    Example: ok remove ./dist --permanent")
	fmt.Println()
	fmt.Println("  ok docker")
	fmt.Println("    Launches an interactive UI to manage Docker containers.")
	fmt.Println()
	fmt.Println("  ok kill --port <port>")
	fmt.Println("    macOS: Finds processes listening on the TCP port, lists them, and asks for confirmation.")
	fmt.Println("    Confirmation prompt defaults to 'Y' on Enter. Kills with SIGKILL (-9) when confirmed.")
	fmt.Println("    Example: ok kill --port 3000")
	fmt.Println()

	color.Yellow("Config:")
	fmt.Println("  Defaults live at ~/.ok/config.yaml")
	fmt.Println()

	color.Yellow("Tips:")
	fmt.Println("  • Use --verbose to see success messages.")
	fmt.Println("  • Run 'ok <command> --help' for command-specific flags.")
}
