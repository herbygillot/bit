package scripts

import (
	"fmt"
	"github.com/chriswalz/complete/v2"
	"github.com/chriswalz/complete/v2/predict"
	"io/ioutil"
	"os"
	"strings"
)

func parseZshAutocompleteOutput(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		return
	}
	s := string(b[:])
	//log.Print(s)
	sg := map[string]*complete.Command{}
	s = strings.ReplaceAll(s, "#######", "#####")
	s = strings.ReplaceAll(s, "######", "#####")
	parts := strings.Split(s, "#####")
	for i := 0; i < len(parts); i += 2 {
		argPart := parts[i]
		flagPart := parts[i+1]
		argPart = strings.TrimSpace(argPart)
		//fmt.Println(argPart)
		lines := strings.Split(argPart, "\n")
		command := strings.TrimSpace(strings.ReplaceAll(lines[0], "$ ", ""))
		command = strings.ReplaceAll(command, " --", "")
		command = strings.ReplaceAll(command, " -", "")
		//fmt.Println(command)
		var subsubs = map[string]*complete.Command{}
		if len(lines) > 1 {
			for _, subb := range lines[1:] {
				if strings.Contains(subb, "%backup%") || strings.Contains(subb, "Applications/") {
					break
				}
				subb = strings.TrimSpace(subb)
				ss := strings.Split(subb, "--")
				desc := strings.TrimSpace(ss[len(ss)-1])
				sub := strings.Fields(strings.TrimSpace(ss[0]))[0]
				subsubs[sub] = &complete.Command{
					Description: desc,
				}
			}
		}
		sub := strings.ReplaceAll(command, "git ", "")

		flagPart = strings.TrimSpace(flagPart)
		lines = strings.Split(flagPart, "\n")

		var flags = map[string]complete.Predictor{}
		if len(lines) > 1 {
			for _, subb := range lines[1:] {
				if strings.Contains(subb, "%backup%") || strings.Contains(subb, "Applications/") {
					break
				}
				subb = strings.TrimSpace(subb)
				ss := strings.Split(subb[2:], "--")
				//desc := strings.TrimSpace(ss[len(ss)-1])
				flagName := strings.TrimSpace(ss[0])
				flags[flagName] = predict.Nothing
			}
		}

		sg[sub] = &complete.Command{
			Description: "",
			Sub:         subsubs,
			Flags:       flags,
			Args:        nil,
		}
	}
	//fmt.Println("hello")
	bittree := &complete.Command{
		Description: "",
		Sub:         sg,
		Flags:       nil,
		Args:        nil,
	}
	codestring := printSuggestionTreeCode(bittree)
	f, err := os.Create("cmd/code_generated_src.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(codestring)

}

func printSuggestionTreeCode(suggestionTree *complete.Command) string {
	flags := ""
	flagsStruct := ""
	if suggestionTree.Flags != nil {
		for flagName := range suggestionTree.Flags {
			flags += "\"" + flagName + "\": predict.Nothing,\n"
		}
		flagsStruct = `
Flags: map[string]complete.Predictor{
` + flags + `
},
`
	}
	subs := ""
	if suggestionTree.Sub != nil {
		for key, val := range suggestionTree.Sub {
			subs += "\"" + key + "\": " + printSuggestionTreeCode(val)

		}
	}
	return `
&complete.Command{ Description: "", Sub: map[string]*complete.Command{` + subs + `}, ` + flagsStruct + `},
`

}
