package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/spf13/cobra"
	"log"
)

var findCmd = &cobra.Command{
	Use:   "find <search-string>",
	Short: "Secretsmanager utilities",
	Long: `Argument:
<search-string>   Keyword to search for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		searchString := args[0]
		cfg := loadConfiguration()
		client := secretsmanager.NewFromConfig(cfg)

		paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{
					Key:    types.FilterNameStringTypeAll,
					Values: []string{searchString},
				},
			},
			IncludePlannedDeletion: nil,
		})

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatal(err)
			}

			for _, object := range output.SecretList {
				log.Printf("key=%s description=%s", aws.ToString(object.Name), aws.ToString(object.Description))
			}
		}
	},
}

func init() {
	findCmd.Flags().BoolP("include-deleted", "d", false, "Include secrets marked for deletion")
}
