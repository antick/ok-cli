package cmd

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

type processInfo struct {
	PID     int
	Command string
	User    string
	Name    string // NAME column from lsof (e.g., *:3000 or 127.0.0.1:3000)
}

// HandleKill implements `ok kill --port <port>`
func HandleKill(cmd *cobra.Command, args []string) {
    port, _ := cmd.Flags().GetInt("port")
    if port == 0 {
        color.Red("Error: No port provided. Use --port <port>.")
        // Show full help so the user can see how to use this command
        _ = cmd.Root().Help()
        return
    }

	verbose, _ := cmd.Flags().GetBool("verbose")

	procs, err := findProcessesOnPort(port)
	if err != nil {
		color.Red("Error finding processes on port %d: %v", port, err)
		return
	}
	if len(procs) == 0 {
		color.Yellow("No processes found listening on port %d", port)
		return
	}
	color.Cyan("Found %d process(es) using port %d:", len(procs), port)
	printProcessTable(procs)

	fmt.Print("Proceed to kill them? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" && (strings.EqualFold(input, "n") || strings.EqualFold(input, "no")) {
		color.Yellow("Aborted.")
		return
	}

	var failed []int
	for _, p := range procs {
		if err := syscall.Kill(p.PID, syscall.SIGKILL); err != nil {
			failed = append(failed, p.PID)
		} else if verbose {
			color.Green("Killed PID %d (%s)", p.PID, p.Command)
		}
	}

	if len(failed) == 0 {
		color.Green("Successfully freed port %d", port)
	} else {
		color.Yellow("Some processes could not be killed: %v", failed)
	}
}

func findProcessesOnPort(port int) ([]processInfo, error) {
    c := exec.Command("lsof", "-nP", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN")
    var out bytes.Buffer
    c.Stdout = &out
    c.Stderr = &out
    _ = c.Run() // lsof may return non-zero when nothing is found; parse output regardless

	_lines := strings.Split(out.String(), "\n")
	var res []processInfo
	for i, line := range _lines {
		line = strings.TrimSpace(line)
		if i == 0 && strings.HasPrefix(strings.ToUpper(line), "COMMAND") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		// NAME is typically the last field, with optional trailing "(LISTEN)" token
		name := fields[len(fields)-1]
		if name == "(LISTEN)" && len(fields) >= 2 {
			name = fields[len(fields)-2]
		}
		res = append(res, processInfo{
			PID:     pid,
			Command: fields[0],
			User:    fields[2],
			Name:    name,
		})
	}

	// Deduplicate by PID
	seen := map[int]bool{}
	uniq := make([]processInfo, 0, len(res))
	for _, p := range res {
		if !seen[p.PID] {
			seen[p.PID] = true
			uniq = append(uniq, p)
		}
	}
	return uniq, nil
}

func printProcessTable(procs []processInfo) {
	cmdWidth := len("COMMAND")
	userWidth := len("USER")
	for _, p := range procs {
		if len(p.Command) > cmdWidth {
			cmdWidth = len(p.Command)
		}
		if len(p.User) > userWidth {
			userWidth = len(p.User)
		}
	}
	header := fmt.Sprintf("%-*s  %-*s  %-5s  %s", cmdWidth, "COMMAND", userWidth, "USER", "PID", "NAME")
	color.Yellow(header)
	for _, p := range procs {
		fmt.Printf("%-*s  %-*s  %-5d  %s\n", cmdWidth, p.Command, userWidth, p.User, p.PID, p.Name)
	}
}
