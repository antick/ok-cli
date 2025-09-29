package main

import (
    "os"

    "github.com/fatih/color"
    "github.com/spf13/cobra"

    "ok/cmd"
    "ok/config"
)

var cfg config.Config

func main() {
    var err error
    cfg, err = config.LoadConfig()
    if err != nil {
        color.Red("Error loading config: %v", err)
        os.Exit(1)
    }

    var rootCmd = &cobra.Command{
        Use:   "ok",
        Long:  `OK CLI is a versatile command-line tool that provides various utilities for file operations, building Go programs, and managing Docker containers.`,
        CompletionOptions: cobra.CompletionOptions{
            DisableDefaultCmd: true,
        },
    }

    rootCmd.PersistentFlags().BoolVarP(&cfg.VerboseOutput, "verbose", "v", cfg.VerboseOutput, "verbose output")

    rootCmd.AddCommand(
        createCopyCommand(),
        createBuildCommand(),
        createMoveCommand(),
        createRemoveCommand(),
        createDockerCommand(),
        createKillCommand(),
        createVersionCommand(),
    )

    // Use our custom help function for consistent, detailed help output
    rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
        cmd.HandleHelp(c, args)
    })

    if err := rootCmd.Execute(); err != nil {
        color.Red("Error: %v", err)
        os.Exit(1)
    }
}

func createVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print the version number of OK CLI",
        Long:  `All software has versions. This is OK CLI's version.`,
        Run:   cmd.HandleVersion,
    }
}

func createCopyCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "copy <source> [to] <destination>",
        Short: "Copy files or directories",
        Long:  `Copy files or directories from source to destination.`,
        Run:   cmd.HandleCopy,
    }
    cmd.Flags().StringVarP(&cfg.DefaultDestination, "destination", "d", cfg.DefaultDestination, "default destination")
    return cmd
}

func createBuildCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "build <input_file> [as/to] <output_file>",
        Short: "Build Go programs",
        Long:  `Build Go programs from source files.`,
        Run:   cmd.HandleBuild,
    }
    cmd.Flags().StringVarP(&cfg.DefaultBuildOutput, "output", "o", cfg.DefaultBuildOutput, "default build output name")
    return cmd
}

func createMoveCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "move <source> [to] <destination>",
        Short: "Move files or directories",
        Long:  `Move files or directories from source to destination.`,
        Run:   cmd.HandleMove,
    }
    cmd.Flags().StringVarP(&cfg.DefaultDestination, "destination", "d", cfg.DefaultDestination, "default destination")
    return cmd
}

func createRemoveCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "remove <file_or_directory>",
        Short: "Remove files or directories",
        Long:  `Remove files or directories, moving them to trash by default.`,
        Run:   cmd.HandleRemove,
    }
    cmd.Flags().BoolP("permanent", "p", false, "permanently delete instead of moving to trash")
    return cmd
}

func createDockerCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "docker",
        Short: "Manage Docker containers",
        Run:   cmd.HandleDocker,
    }
}

func createKillCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "kill [--port] <port>",
        Short: "Kill processes listening on a TCP port",
        Long:  `Finds processes listening on the given TCP port, shows them, and prompts for confirmation before killing. You can specify the port either as a flag (--port 3000) or as a positional argument (3000).`,
        Run:   cmd.HandleKill,
    }
    cmd.Flags().IntP("port", "p", 0, "TCP port to free (required)")
    return cmd
}
