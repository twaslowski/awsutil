package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var findCmd = &cobra.Command{
	Use:   "find <search-string>",
	Short: "Fuzzy-find secrets",
	Long: `Argument:
<search-string>   Keyword to search for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:  executeFind(),
}

func executeFind() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		searchString := args[0]
		includeDeleted, _ := cmd.Flags().GetBool("include-deleted")
		cfg := loadConfiguration()
		client := secretsmanager.NewFromConfig(cfg)

		paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{
					Key:    types.FilterNameStringTypeAll,
					Values: []string{searchString},
				},
			},
			IncludePlannedDeletion: &includeDeleted,
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Description"})
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatal(err)
			}

			for _, object := range output.SecretList {
				table.Append([]string{aws.ToString(object.Name), aws.ToString(object.Description)})
			}
		}
		table.Render()
	}
}

func init() {
	findCmd.Flags().BoolP("include-deleted", "d", false, "Include secrets marked for deletion")
}
