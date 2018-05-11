package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/tcnksm/go-gitconfig"
	"github.com/techjacker/diffence"
)

const (
	defaultHooksPath = ".git/hooks/"
	rulesFile        = "precommit-rules.json"
	ignoreFile       = ".precommit-ignore"
)

var (
	err       error
	rules     *[]diffence.Rule
	hooksPath string
)

func main() {
	hooksPath, err = gitconfig.Global("core.hooksPath")
	if hooksPath == "" {
		hooksPath = defaultHooksPath
	}

	stagedChanges, err := exec.Command("/usr/bin/git", "diff", "--staged").Output()
	if err != nil {
		log.Fatalf("Cannot get staged changes\n%s", err)
		return
	}

	rules, err = diffence.LoadRulesJSON(hooksPath + rulesFile)
	if err != nil {
		rules, err = diffence.LoadDefaultRules()
		if err != nil {
			log.Fatalf("Cannot load default rules\n%s", err)
			return
		}
	}

	diff := diffence.DiffChecker{
		Rules:   rules,
		Ignorer: diffence.NewIgnorerFromFile(ignoreFile),
	}

	res, err := diff.Check(bytes.NewReader(stagedChanges))
	if err != nil {
		log.Fatalf("Error reading diff\n%s\n", err)
		return
	}

	matches := res.Matches()
	if matches < 1 {
		os.Exit(0)
	}

	i := 1
	fmt.Fprintf(os.Stderr, "git-rid-of-keys\n\ncurrent commit contains %d offenses!\nremove or add the path to `.precommit-ignore` file right meow!\n", matches)
	for diffKey, rule := range res.MatchedRules {
		fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------\n")
		fmt.Fprintf(os.Stderr, "offense: #%d\n", i)
		commit, filename := diffence.SplitDiffHashKey(diffKey)
		if commit != "" {
			fmt.Fprintf(os.Stderr, "commit: %s\n", commit)
		}
		fmt.Fprintf(os.Stderr, "file: %s\n", filename)
		fmt.Fprintf(os.Stderr, "reason: %#v\n", rule[0].Caption)
		i++
	}
	os.Exit(1)
}
