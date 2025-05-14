package logs

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "logs",
	Short: "logging analysis",
}

func init() {
	Cmd.AddCommand(flowLogsCmd)
}
