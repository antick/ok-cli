package cmd

import (
	"github.com/antick/ok/docker"

	"github.com/spf13/cobra"
)

func HandleDocker(cmd *cobra.Command, args []string) {
	docker.RunDockerUI()
}
