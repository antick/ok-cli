package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunDockerUI() {
	app := tview.NewApplication()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerList := tview.NewTable().SetSelectable(true, false).SetBorders(true)
	containerList.SetTitle("Containers").SetBorder(true)

	statsView := tview.NewTextView().SetDynamicColors(true)
	statsView.SetTitle("Container Stats").SetBorder(true)

	logView := tview.NewTextView().SetDynamicColors(true)
	logView.SetTitle("Container Logs").SetBorder(true)

	helpBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[::b]=[::-]:refresh  [::b]i[::-]:info  [::b]l[::-]:shell  [::b][RETURN][::-]:logs  [::b]r[::-]:restart  [::b]s[::-]:stop  [::b]:q[::-]:quit")

	dashboardPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(containerList, 0, 3, true).
			AddItem(statsView, 0, 1, false),
			0, 1, true).
		AddItem(logView, 0, 1, false).
		AddItem(helpBar, 1, 0, false)

	overlay := tview.NewBox().
		SetBorder(true).
		SetTitle("Command Input").
		SetTitleAlign(tview.AlignLeft)

	commandInput := tview.NewInputField().
		SetLabel(":").
		SetFieldWidth(2).
		SetFieldBackgroundColor(tcell.ColorBlack)

	overlayFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().AddItem(commandInput, 0, 1, true), 1, 1, false).
			AddItem(nil, 0, 1, false),
			3, 1, true,
		).
		AddItem(nil, 0, 1, false)

	overlay.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		overlayFlex.SetRect(x, y, width, height)
		overlayFlex.Draw(screen)
		return x, y, width, height
	})

	overlayPage := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().AddItem(commandInput, 0, 0, false), 1, 1, false).
			AddItem(overlay, 3, 1, false).
			AddItem(nil, 0, 1, false), 40, 1, true).
		AddItem(nil, 0, 1, false)

	pages := tview.NewPages().
		AddPage("main", dashboardPage, true, true).
		AddPage("input", overlayPage, true, false)

	commandInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := commandInput.GetText()
			if command == "q" {
				app.Stop()
			}

			commandInput.SetText("")
			pages.HidePage("input")
			app.SetFocus(containerList)
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.HidePage("input")
			app.SetFocus(containerList)
			return nil
		}
		if event.Rune() == ':' {
			pages.ShowPage("input")
			app.SetFocus(commandInput)
			return nil
		}
		return event
	})

	updateContainers := func() {
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return
		}

		containerList.Clear()
		containerList.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow))
		containerList.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow))
		containerList.SetCell(0, 2, tview.NewTableCell("Image").SetTextColor(tcell.ColorYellow))
		containerList.SetCell(0, 3, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow))

		for i, container := range containers {
			containerList.SetCell(i+1, 0, tview.NewTableCell(container.ID[:12]))
			containerList.SetCell(i+1, 1, tview.NewTableCell(container.Names[0][1:]))
			containerList.SetCell(i+1, 2, tview.NewTableCell(container.Image))
			containerList.SetCell(i+1, 3, tview.NewTableCell(container.Status))
		}
	}

	updateStats := func(containerID string) {
		stats, err := cli.ContainerStats(context.Background(), containerID, false)
		if err != nil {
			return
		}
		defer stats.Body.Close()

		var statsJSON types.StatsJSON
		if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
			return
		}

		cpuPercent := calculateCPUPercentUnix(statsJSON.CPUStats, statsJSON.PreCPUStats)
		memoryUsage := float64(statsJSON.MemoryStats.Usage) / 1024 / 1024
		memoryLimit := float64(statsJSON.MemoryStats.Limit) / 1024 / 1024

		statsView.Clear()
		fmt.Fprintf(statsView, "CPU: %.2f%%\nMemory: %.2f / %.2f MB\n", cpuPercent, memoryUsage, memoryLimit)
	}

	updateLogs := func(containerID string) {
		logs, err := cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "10",
		})
		if err != nil {
			return
		}
		defer logs.Close()

		logView.Clear()
		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			fmt.Fprintln(logView, scanner.Text())
		}
	}

	containerList.SetSelectedFunc(func(row int, column int) {
		if row > 0 {
			containerID := containerList.GetCell(row, 0).Text
			updateStats(containerID)
			updateLogs(containerID)
		}
	})

	app.SetFocus(containerList)
	app.SetRoot(pages, true).EnableMouse(true)

	go func() {
		for {
			app.QueueUpdateDraw(func() {
				updateContainers()
			})
			time.Sleep(5 * time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func calculateCPUPercentUnix(v container.CPUStats, pre container.CPUStats) float64 {
	cpuPercent := 0.0
	cpuDelta := float64(v.CPUUsage.TotalUsage) - float64(pre.CPUUsage.TotalUsage)
	systemDelta := float64(v.SystemUsage) - float64(pre.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}
