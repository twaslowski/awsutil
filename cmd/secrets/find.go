package secrets

import (
	util "awsutil/pkg"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var findCmd = &cobra.Command{
	Use:   "find <search-string>",
	Short: "Fuzzy-find secrets",
	Long: `Argument:
<search-string>   Keyword to search for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: executeFind(),
}

func executeFind() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		searchString := args[0]

		region, err := cmd.Flags().GetString("region")
		if err != nil {
			return err
		}

		includeDeleted, err := cmd.Flags().GetBool("include-deleted")
		if err != nil {
			return err
		}

		showArn, err := cmd.Flags().GetBool("show-arn")
		if err != nil {
			return err
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return err
		}

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
		table.Header(getHeaders(showArn))

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				return err
			}

			for _, object := range output.SecretList {
				err = table.Append(getTableEntry(object, showArn))
				if err != nil {
					return err
				}
			}
		}
		return table.Render()
	}
}

func getTableEntry(entry types.SecretListEntry, showArn bool) []string {
	identifier := aws.ToString(entry.Name)
	if showArn {
		identifier = aws.ToString(entry.ARN)
	}
	row := []string{identifier, util.TruncateString(aws.ToString(entry.Description), 50)}
	return row
}

func getHeaders(showArn bool) []string {
	headers := []string{"Description"}
	if showArn {
		headers = append([]string{"ARN"}, headers...)
	} else {
		headers = append([]string{"Name"}, headers...)
	}
	return headers
}

func init() {
	findCmd.Flags().Bool("include-deleted", false, "Include secrets marked for deletion")
	findCmd.Flags().Bool("show-arn", false, "Include ARN in the result table")
}
