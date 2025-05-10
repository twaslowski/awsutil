package secrets

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

type CloudTrailEvent struct {
	EventSource     string
	EventName       string
	UserIdentity    UserIdentity
	EventTime       time.Time
	SourceIpAddress string
}

type UserIdentity struct {
	Type        string
	PrincipalId string
	Arn         string
	AccountId   string
}

var lastAccessedCmd = &cobra.Command{
	Use:   "last-accessed <secret-arn>",
	Short: "Analyze access to secrets",
	Long: `Argument:
<secret-arn>   Secret to check access for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:  executeLastAccess(),
}

func executeLastAccess() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		secret := args[0]
		cfg := loadConfiguration()
		client := cloudtrail.NewFromConfig(cfg)
		startDate := time.Now().AddDate(0, 0, -7)

		paginator := cloudtrail.NewLookupEventsPaginator(client, &cloudtrail.LookupEventsInput{
			StartTime: &startDate,
			EndTime:   nil,
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   "ResourceName",
					AttributeValue: &secret,
				},
			},
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Time", "Event", "Principal"})
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatal(err)
			}

			for _, event := range output.Events {
				var cloudTrailEvent CloudTrailEvent
				err = json.Unmarshal([]byte(aws.ToString(event.CloudTrailEvent)), &cloudTrailEvent)

				if err != nil {
					log.Fatal(err)
				}

				table.Append([]string{cloudTrailEvent.EventTime.String(),
					cloudTrailEvent.EventName,
					cloudTrailEvent.UserIdentity.Arn,
				})
			}
		}
		table.Render()
	}
}
