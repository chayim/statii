package plugins

import (
	"context"
	"strings"

	"github.com/google/go-github/v48/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

var GITHUB_ISSUE_CONFIG string = "github_issue"

type GitHubIssueConfig struct {
	Token        string
	Repositories []string
	States       []string
}

func NewGitHubIssueConfig(m interface{}) *GitHubIssueConfig {
	g := GitHubIssueConfig{}
	g.Token = m["Token"]
	g.Repositories = m["Repositories"]
	g.States = m["States"]
	return &g
}

func (g *GitHubIssueConfig) Execute(ctx context.Context) {
	if g.Token == "" || len(g.Repositories) == 0 {
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for _, repo := range g.Repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			log.Warnf("%s is an invalid repository, skipping.", repo)
			continue
		}

		opts := github.IssueListByRepoOptions{State: "all", Sort: "updated", Direction: "desc"}
		issues, _, err := client.Issues.ListByRepo(ctx, parts[0], parts[1], &opts)
		if err != nil {
			log.Errorf("error on %s: %v", repo, err)
			continue
		}

		for _, p := range issues {
			if slices.Contains(g.States, *p.State) {
				// TODO notify
			}
		}

	}

}
