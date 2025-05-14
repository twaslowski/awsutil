package logs

import "github.com/spf13/cobra"

var flowLogsCmd = &cobra.Command{
	Use:   "find <search-string>",
	Short: "Fuzzy-find secrets",
	Long: `Argument:
<search-string>   Keyword to search for`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: executeFlowLogs(),
}

func executeFlowLogs() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
