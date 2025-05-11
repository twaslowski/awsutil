package cmd

import (
	"awsutil/cmd/trace"
	"os"

	"awsutil/cmd/secrets"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsutil",
	Short: "An opinionated AWS utility for everyday tasks",
	Long: `awsutil provides utilities around several services relevant to the author's everyday work.
It includes services such as CloudWatch, Secretsmanager, SQS and EKS.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(secrets.Cmd)
	rootCmd.AddCommand(trace.Cmd)
}
