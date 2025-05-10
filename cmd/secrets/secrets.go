package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"log"
)

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "Secretsmanager utilities",
}

func loadConfiguration() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func init() {
	Cmd.AddCommand(findCmd)
	Cmd.AddCommand(lastAccessedCmd)
}
