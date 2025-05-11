package secrets

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "Secretsmanager utilities",
}

func init() {
	Cmd.AddCommand(findCmd)
}
