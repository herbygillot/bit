package cmd

import (
	"github.com/chriswalz/complete/v2"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Check out a pull request from Github (requires GH CLI)",
	Long: `bit pr
bit pr`,
	Run: func(cmd *cobra.Command, args []string) {
		suggestionTree := &complete.Command{
			Sub: map[string]*complete.Command{
				"pr": {
					Description: "Check out a pull request from Github (requires GH CLI)",
					Args:        complete.PredictFunc(lazyLoad(GitHubPRSuggestions)),
				},
			},
		}
		runPr(suggestionTree)
	},
	Args: cobra.NoArgs,
}

func init() {
	BitCmd.AddCommand(prCmd)
}

func runPr(suggestionMap *complete.Command) {
	branchName := SuggestionPrompt("> bit pr ", specificCommandCompleter("pr", suggestionMap))

	split := strings.Split(branchName, "#")
	prNumber, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		log.Debug().Err(err)
		return
	}
	checkoutPullRequest(prNumber)
	return
}
