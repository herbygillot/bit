package cmd

import (
	"fmt"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/thoas/go-funk"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// ShellCmd represents the base command when called without any subcommands
var ShellCmd = &cobra.Command{
	Use:   "bit",
	Short: "Bit is a Git CLI that predicts what you want to do",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		suggestionTree, bitCmdMap := CreateSuggestionMap(cmd)

		resp := SuggestionPrompt("> bit ", shellCommandCompleter(suggestionTree))
		subCommand := resp
		if subCommand == "" {
			return
		}
		if strings.Index(resp, " ") > 0 {
			subCommand = subCommand[0:strings.Index(resp, " ")]
		}
		parsedArgs, err := parseCommandLine(resp)
		if err != nil {
			log.Debug().Err(err)
			return
		}
		if bitCmdMap[subCommand] == nil {
			yes := HijackGitCommandOccurred(parsedArgs, suggestionTree, cmd.Version)
			if yes {
				return
			}
			RunGitCommandWithArgs(parsedArgs)
			return
		}

		cmd.SetArgs(parsedArgs)
		cmd.Execute()
	},
}

func init() {
	ShellCmd.PersistentFlags().Bool("debug", false, "Print debugging information")
}

func CreateSuggestionMap(cmd *cobra.Command) (*BitCommand, map[string]*cobra.Command) {
	start := time.Now()
	_, bitCmdMap := AllBitSubCommands(cmd)
	log.Debug().Msg((time.Now().Sub(start)).String())
	start = time.Now()
	allBitCmds := AllBitAndGitSubCommands(cmd)
	log.Debug().Msg((time.Now().Sub(start)).String())
	//commonCommands := CobraCommandToSuggestions(CommonCommandsList())
	start = time.Now()
	branchListSuggestions := BranchListSuggestions()
	log.Debug().Msg((time.Now().Sub(start)).String())
	start = time.Now()
	_ = CobraCommandToSuggestions(allBitCmds)
	log.Debug().Msg((time.Now().Sub(start)).String())
	start = time.Now()
	//gitAddSuggestions := GitAddSuggestions()
	log.Debug().Msg((time.Now().Sub(start)).String())
	start = time.Now()
	//gitResetSuggestions := GitResetSuggestions()
	log.Debug().Msg((time.Now().Sub(start)).String())
	start = time.Now()
	//gitmoji := GitmojiSuggestions()
	log.Debug().Msg((time.Now().Sub(start)).String())

	branchListText := funk.Map(branchListSuggestions, func(s prompt.Suggest) string {
		return s.Text
	}).([]string)

	suggestionTree := &BitCommand{
		//Args: combraCommandSuggestions,
		Flags: map[string]complete.Predictor{
			"--version": predict.Nothing,
		},
		Sub: map[string]*BitCommand{
			"add":          {Description: "Add file contents to the index"},
			"am":           {Description: "Apply a series of patches from a mailbox"},
			"archive":      {Description: "Create an archive of files from a named tree"},
			"branch":       {Description: "List, create, or delete branches"},
			"bisect":       {Description: "Use binary search to find the commit that introduced a bug"},
			"bundle":       {Description: "Move objects and refs by archive"},
			"commit":       {Description: "Record changes to the repository"},
			"clone":        {Description: "Clone a repository into a new directory"},
			"checkout":     {Description: "Switch branches or restore working tree files", Args: predict.Set(branchListText)},
			"fetch":        {Description: "Download objects and refs from another repository"},
			"diff":         {Description: "Show changes between commits, commit and working tree, etc"},
			"cherry-pick":  {Description: "Apply the changes introduced by some existing commits"},
			"citool":       {Description: "Graphical alternative to git-commit"},
			"clean":        {Description: "Remove untracked files from the working tree"},
			"describe":     {Description: "Give an object a human readable name based on an available ref"},
			"format-patch": {Description: "Prepare patches for e-mail submission"},
			"gc":           {Description: "Cleanup unnecessary files and optimize the local repository"},
			"gitk":         {Description: "The Git repository browser"},
			"grep":         {Description: "Print lines matching a pattern"},
			"gui":          {Description: "A portable graphical interface to Git"},
			"init":         {Description: "Create an empty Git repository or reinitialize an existing one"},
			"log":          {Description: "Show commit logs", Args: predict.Set(branchListText)},
			"merge":        {Description: "Join two or more development histories together", Args: predict.Set(branchListText)},
			"mv":           {Description: "Move or rename a file, a directory, or a symlink"},
			"notes":        {Description: "Add or inspect object notes"},
			"pull":         {Description: "Fetch from and integrate with another repository or a local branch"},
			"push":         {Description: "Update remote refs along with associated objects"},
			"range-diff":   {Description: "Compare two commit ranges (e.g. two versions of a branch)"},
			"rebase":       {Description: "Reapply commits on top of another base tip", Args: predict.Set(branchListText)},
			"release": {
				Description: "Commit unstaged changes, bump minor tag, push",
				Args:        predict.Set{"bump", "<version>"},
			},
			"pr": {
				Description: "Check out a pull request from Github (requires GH CLI)",
				Args:        complete.PredictFunc(test123),
			},
			"reset": {Description: "Reset current HEAD to the specified state",
				Flags: map[string]complete.Predictor{
					"--soft": predict.Nothing,
				},
				Args: predict.Set{"HEAD~1", "gc"}},
			"restore":  {Description: "Restore working tree files"},
			"revert":   {Description: "Revert some existing commits"},
			"rm":       {Description: "Remove files from the working tree and from the index"},
			"show":     {Description: "Show various types of objects"},
			"stash":    {Description: "Stash the changes in a dirty working directory away"},
			"shortlog": {Description: "Summarize 'git log' output"},
			//"status": {Description: "Show the working tree status"}, NOTE: mved common suggestions to keep close to top
			"submodule":     {Description: "Initialize, update or inspect submodules"},
			"switch":        {Description: "Switch branches", Args: predict.Set(branchListText)},
			"tag":           {Description: "Create, list, delete or verify a tag object signed with GPG"},
			"worktree":      {Description: "Manage multiple working trees"},
			"config":        {Description: "Get and set repository or global options"},
			"fast-import":   {Description: "Backend for fast Git data importers"},
			"filter-branch": {Description: "Rewrite branches"},
			"mergetool":     {Description: "Run merge conflict resolution tools to resolve merge conflicts"},
			"pack-refs":     {Description: "Pack heads and tags for efficient repository access"},
			"prune":         {Description: "Prune all unreachable objects from the object database"},
			"reflog":        {Description: "Manage reflog information"},
			"remote": {Description: "Manage set of tracked repositories",
				Sub: map[string]*BitCommand{
					"rename":   {},
					"remove":   {},
					"set-head": {},
				}},
			"repack":          {Description: "Pack unpacked objects in a repository"},
			"replace":         {Description: "Create, list, delete refs to replace objects"},
			"annotate":        {Description: "Annotate file lines with commit information"},
			"blame":           {Description: "Show what revision and author last modified each line of a file"},
			"count-objects":   {Description: "Count unpacked number of objects and their disk consumption"},
			"difftool":        {Description: "Show changes using common diff tools"},
			"fsck":            {Description: "Verifies the connectivity and validity of the objects in the database"},
			"gitweb":          {Description: "Git web interface (web frontend to Git repositories)"},
			"help":            {Description: "Display help information about Git"},
			"instaweb":        {Description: "Instantly browse your working repository in gitweb"},
			"merge-tree":      {Description: "Show three-way merge without touching index"},
			"rerere":          {Description: "Reuse recorded resolution of conflicted merges"},
			"show-branch":     {Description: "Show branches and their commits"},
			"verify-commit":   {Description: "Check the GPG signature of commits"},
			"verify-tag":      {Description: "Check the GPG signature of tags"},
			"whatchanged":     {Description: "Show logs with difference each commit introduces"},
			"archimport":      {Description: "Import a GNU Arch repository into Git"},
			"cvsexportcommit": {Description: "Export a single commit to a CVS checkout"},
			"cvsimport":       {Description: "Salvage your data out of another SCM people love to hate"},
			"cvsserver":       {Description: "A CVS server emulator for Git"},
			"imap-send":       {Description: "Send a collection of patches from stdin to an IMAP folder"},
			"p4":              {Description: "Import from and submit to Perforce repositories"},
			"fast-export":     {Description: "Git data exporter"},
			"version":         {Description: "Print bit and git version"},
		},
	}

	// command
	// flags
	// commands
	// value

	//completerSuggestionMap := map[string]func() []prompt.Suggest{
	//	"":         memoize([]prompt.Suggest{}),
	//	"shell":    memoize(combraCommandSuggestions),
	//	"checkout": memoize(branchListSuggestions),
	//	"switch":   memoize(branchListSuggestions),
	//	"co":       memoize(branchListSuggestions),
	//	"merge":    memoize(branchListSuggestions),
	//	"rebase":   memoize(branchListSuggestions),
	//	"log":      memoize(branchListSuggestions),
	//	"add":      memoize(gitAddSuggestions),
	//	"release": memoize([]prompt.Suggest{
	//		{Text: "bump", Description: "Increment SemVer from tags and release e.g. if latest is v0.1.2 it's bumped to v0.1.3 "},
	//		{Text: "<version>", Description: "Name of release version e.g. v0.1.2"},
	//	}),
	//	"reset":   memoize(gitResetSuggestions),
	//	"pr":      lazyLoad(GitHubPRSuggestions),
	//	"gitmoji": memoize(gitmoji),
	//	"save":    memoize(gitmoji),
	//	//"_any": commonCommands,
	//}
	return suggestionTree, bitCmdMap
}

func test123(prefix string) []string {
	return []string{"example-pr"}
}

// Execute adds all child commands to the shell command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the ShellCmd.
func Execute() {
	if err := ShellCmd.Execute(); err != nil {
		log.Info().Err(err)
		os.Exit(1)
	}
}

func shellCommandCompleter(suggestionTree *BitCommand) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return promptCompleter(suggestionTree, d.Text)
	}
}

func branchCommandCompleter(suggestionMap *BitCommand) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return promptCompleter(suggestionMap, "checkout "+d.Text)
	}
}

func specificCommandCompleter(subCmd string, suggestionMap *BitCommand) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return promptCompleter(suggestionMap, subCmd+" "+d.Text)
	}
}

func promptCompleter(suggestionTree *BitCommand, text string) []prompt.Suggest {
	split := strings.Split("bit "+text, " ")
	currCmdNode := suggestionTree

	// the last token is the query/prefix
	query := split[len(split)-1]
	var suggestionsFunc complete.Predictor = nil
	// traverse tree to find the last specified sub command from input
	for _, token := range split {
		keys := make([]string, 0, len(currCmdNode.Sub))
		for k := range currCmdNode.Sub {
			keys = append(keys, k)
		}
		var args []string
		if currCmdNode.Args != nil {
			args = currCmdNode.Args.Predict("")
		}
		if currCmdNode.Flags != nil && strings.HasPrefix(query, "-") {
			keys := make([]string, 0, len(currCmdNode.Flags))
			for k := range currCmdNode.Flags {
				keys = append(keys, k)
			}
			args = append(args, keys...)
		}
		suggestionsFunc = predict.Set(append(args, keys...))
		if currCmdNode == nil {
			suggestionsFunc = predict.Set{"command-not-configured-notify-developer"}
			break
		}
		if currCmdNode.Sub[token] != nil {
			currCmdNode = currCmdNode.Sub[token]
			continue
		}
	}

	suggestions := suggestionsFunc.Predict("")
	var sugg []prompt.Suggest
	for _, suggestion := range suggestions {
		sugg = append(sugg, prompt.Suggest{Text: suggestion})
	}

	//if len(prev) == len(text) {
	//	suggestions = suggestionTree["shell"]()
	//	return prompt.FilterHasPrefix(suggestions, prev, true)
	//}
	//curr := filterFlags[1] // in git commit -m "hello"  "hello" is curr
	//if strings.HasPrefix(curr, "--") {
	//	suggestions = FlagSuggestionsForCommand(prev, "--")
	//} else if strings.HasPrefix(curr, "-") {
	//	suggestions = FlagSuggestionsForCommand(prev, "-")
	//} else if suggestionTree[prev] != nil {
	//	suggestions = suggestionTree[prev]()
	//	if isBranchCompletionCommand(prev) || prev == "gitmoji" {
	//		return prompt.FilterContains(suggestions, curr, true)
	//	}
	//}
	return prompt.FilterHasPrefix(sugg, query, true)
}

func RunGitCommandWithArgs(args []string) {
	var err error
	err = RunInTerminalWithColor("git", args)
	if err != nil {
		log.Debug().Msg("Command may not exist: " + err.Error())
	}
	return
}

func HijackGitCommandOccurred(args []string, suggestionMap *BitCommand, version string) bool {
	sub := args[0]
	// handle checkout,switch,co commands as checkout
	// if "-b" flag is not provided and branch does not exist
	// user would be prompted asking whether to create a branch or not
	// expected usage format
	//   bit (checkout|switch|co) [-b] branch-name
	if args[len(args)-1] == "--version" || args[len(args)-1] == "version" {
		fmt.Println("bit version " + version)
		return false
	}
	if sub == "pr" {
		runPr(suggestionMap)
		return true
	}
	if sub == "merge" && len(args) == 1 {
		branchName := SuggestionPrompt("> bit "+sub+" ", specificCommandCompleter("merge", suggestionMap))
		RunInTerminalWithColor("git", []string{"merge", branchName})
		return true
	}
	if isBranchChangeCommand(sub) {
		branchName := ""
		if len(args) < 2 {
			branchName = SuggestionPrompt("> bit "+sub+" ", branchCommandCompleter(suggestionMap))
		} else {
			branchName = strings.TrimSpace(args[len(args)-1])
		}

		if strings.HasPrefix(branchName, "origin/") {
			branchName = branchName[7:]
		}
		args[len(args)-1] = branchName
		var createBranch bool
		if len(args) == 3 && args[len(args)-2] == "-b" {
			createBranch = true
		}
		branchExists := checkoutBranch(branchName)
		if branchExists {
			refreshBranch()
			return true
		}

		if !createBranch && !AskConfirm("Branch does not exist. Do you want to create it?") {
			fmt.Printf("Cancelling...")
			return true
		}

		RunInTerminalWithColor("git", []string{"checkout", "-b", branchName})
		return true
	}
	return false
}

func GetVersion() string {
	return ShellCmd.Version
}

type BitCommand struct {
	Sub         map[string]*BitCommand
	Flags       map[string]complete.Predictor
	Args        complete.Predictor
	Description string
}
