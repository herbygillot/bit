// Package main is complete tool for the go command line
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
)

func Bitcomplete() {
	compLine := os.Getenv("COMP_LINE")
	compPoint := os.Getenv("COMP_POINT")
	doInstall := os.Getenv("COMP_INSTALL") == "1"
	doUninstall := os.Getenv("COMP_UNINSTALL") == "1"

	bitcompletionNotNeeded := compLine == "" && compPoint == "" && !doInstall && !doUninstall
	if bitcompletionNotNeeded {
		return
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	branchCompletion := &BitCommand{
		Args: complete.PredictFunc(func(prefix string) []string {
			branches := BranchListSuggestions()
			completion := make([]string, len(branches))
			for i, v := range branches {
				completion[i] = v.Text
			}
			return completion
		}),
	}

	cmds := AllBitAndGitSubCommands(ShellCmd)
	completionSubCmdMap := map[string]*BitCommand{}
	for _, v := range cmds {
		flagSuggestions := append(FlagSuggestionsForCommand(v.Name(), "--"), FlagSuggestionsForCommand(v.Name(), "-")...)
		flags := funk.Map(flagSuggestions, func(x prompt.Suggest) (string, complete.Predictor) {
			if strings.HasPrefix(x.Text, "--") {
				return x.Text[2:], predict.Nothing
			} else if strings.HasPrefix(x.Text, "-") {
				return x.Text[1:2], predict.Nothing
			} else {
				return "", predict.Nothing
			}
		}).(map[string]complete.Predictor)
		completionSubCmdMap[v.Name()] = &BitCommand{
			Flags: flags,
		}
		if v.Name() == "checkout" || v.Name() == "co" || v.Name() == "switch" || v.Name() == "pull" || v.Name() == "merge" {
			branchCompletion.Flags = flags
			completionSubCmdMap[v.Name()] = branchCompletion
		}
		if v.Name() == "release" {
			completionSubCmdMap[v.Name()].Sub = map[string]*BitCommand{
				"bump": {},
				"test": {},
			}
		}
	}

	gogo := &BitCommand{
		Sub: completionSubCmdMap,
		Flags: map[string]complete.Predictor{
			"version": predict.Nothing,
		},
	}
	//
	//gogo.Complete("bit")
	fmt.Println("fixme", gogo)
}
