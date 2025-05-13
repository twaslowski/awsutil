package eks

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "eks",
	Short: "EKS Utilities",
}

func init() {
	Cmd.AddCommand(selectContextCommand)
}
