package trace

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "trace",
	Short: "Cloudtrail & Diagnostics utilities",
}

func init() {
	Cmd.PersistentFlags().String("start", "7d", "Start date")

	Cmd.AddCommand(accessCmd)
}
