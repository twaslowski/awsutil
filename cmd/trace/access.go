package trace

import (
	util "awsutil/pkg"
	"context"
	"encoding/json"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

var accessCmd = &cobra.Command{
	Use:   "access <resource-arn>",
	Short: "Analyze access on resources",
	Long: `Argument:
<resource-arn>   Resource to check access for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: executeAccessCmd(),
}

func executeAccessCmd() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		secret := args[0]

		region := util.Require(cmd.Flags().GetString("region"))
		start := util.Require(cmd.Flags().GetString("start"))
		startDate, err := parseStartDate(start)
		if err != nil {
			log.Fatalf("Error parsing date: %s", startDate)
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return err
		}

		client := cloudtrail.NewFromConfig(cfg)

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

		table := tablewriter.NewTable(os.Stdout)
		table.Header([]string{"Time", "Event", "Principal"})

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("Failed to retrieve page: %s", err)
			}

			for _, event := range output.Events {
				var cloudTrailEvent CloudTrailEvent
				err = json.Unmarshal([]byte(aws.ToString(event.CloudTrailEvent)), &cloudTrailEvent)
				if err != nil {
					log.Fatalf("Failed to parse JSON: %s", err)
				}

				util.CheckErr(table.Append([]string{
					cloudTrailEvent.EventTime.String(),
					cloudTrailEvent.EventName,
					cloudTrailEvent.UserIdentity.Arn,
				}))
			}
		}

		return table.Render()
	}
}

func parseStartDate(input string) (time.Time, error) {
	startDate, err := dateparse.ParseAny(input)
	if err != nil {
		return calculateStartDateFromDelta(input)
	}
	return startDate, nil
}

func calculateStartDateFromDelta(input string) (time.Time, error) {
	delta, err := util.ParseTimeDelta(input)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().Add(-delta), nil
}
