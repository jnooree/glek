package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	ghClient *github.Client
	ghCtx    context.Context

	version = "master"
	commit  = "none"
	date    = "unknown"

	replacements map[string]bool // Keep tracks which label that has replacement already.
)

// default GitHub labels.
var defaultLabels = [...]string{
	"bug",
	"duplicate",
	"enhancement",
	"help wanted",
	"invalid",
	"question",
	"wontfix",
}

type Label struct {
	Name    string `json:"name"`
	Color   string `json:"color"`
	Replace string `json:"replace,omitempty"`
}

type Gembel struct {
	Labels       []Label  `json:"labels"`
	Repositories []string `json:"repositories"`
}

func main() {
	if len(os.Args) < 2 {
		usage(errors.New("missing <owner/repo>"))
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		usage(errors.New("empty GITHUB_TOKEN in env"))
	}

	ghCtx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ghCtx, ts)
	ghClient = github.NewClient(tc)

	labels, err := getRepoLabels(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	out(labels)
}

// out print the labels into gembel JSON format.
func out(labels []Label) {
	gembel := Gembel{
		Labels:       labels,
		Repositories: []string{"owner/repo"},
	}
	b, err := json.Marshal(gembel)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	json.Indent(&buf, b, "", "\t")
	buf.WriteTo(os.Stdout)
}

// getRepoLabels gets GitHub repository labels.
func getRepoLabels(repoPath string) (labels []Label, err error) {
	parts := strings.Split(repoPath, "/")
	owner, repo := parts[0], parts[1]
	opt := &github.ListOptions{
		PerPage: 100,
	}

	replacements = make(map[string]bool)

	var label Label
	for {
		repoLabels, resp, err := ghClient.Issues.ListLabels(ghCtx, owner, repo, opt)
		if err != nil {
			return labels, err
		}
		for _, ghLabel := range repoLabels {
			label = Label{
				Name:  ghLabel.GetName(),
				Color: ghLabel.GetColor(),
			}
			replaceDefaults(&label)

			labels = append(labels, label)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return labels, err
}

// replaceDefaults replaces default GitHub labels.
func replaceDefaults(label *Label) {
	for _, ghLabel := range defaultLabels {
		if canReplace(label, ghLabel) {
			label.Replace = ghLabel
			break
		}
	}
}

// canReplace checks whether label can replace ghLabel
func canReplace(label *Label, ghLabel string) bool {
	if v, ok := replacements[ghLabel]; v && ok {
		return false
	}

	replacements[ghLabel] = false
	if strings.Contains(strings.ToLower(label.Name), ghLabel) {
		replacements[ghLabel] = true
	}
	return replacements[ghLabel]
}

func getVersion() string {
	return fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
}

func usage(err error) {
	fmt.Printf("Error: %v\n", err)
	fmt.Printf(`
Name:
  glek - export GitHub issue labels into gembel JSON format

Version:
  %s

Usage:
  glek <owner/repo>

  To specifiy GITHUB_TOKEN when running it:

  GITHUB_TOKEN=token glek <owner/repo>
`, getVersion())

	os.Exit(1)
}
