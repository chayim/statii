package plugins

import (
	"context"
	"strconv"
	"strings"

	"github.com/chayim/statii/internal/comms"
	"github.com/google/go-github/v48/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

var GITHUB_PR_CONFIG string = "github_pullrequest"

type GitHubPullRequestConfig struct {
	Token        string   `yaml:"token"`
	Repositories []string `yaml:"repositories"`
	States       []string `yaml:"states,omitempty"`
}

func NewGitHubPullRequestConfig(m interface{}) *GitHubIssueConfig {
	g := GitHubIssueConfig{}
	g.Token = m["Token"]
	g.Repositories = m["Repositories"]
	g.States = m["States"]
	return &g
}

func (g *GitHubPullRequestConfig) Execute(ctx context.Context) {
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

		opts := github.PullRequestListOptions{State: "all", Sort: "updated", Direction: "desc"}
		pulls, _, err := client.PullRequests.List(ctx, parts[0], parts[1], &opts)
		if err != nil {
			log.Errorf("error on %s: %v", repo, err)
			continue
		}

		for _, p := range pulls {
			if slices.Contains(g.States, *p.State) {
				p.GetID()
				msg := comms.NewMessage(strconv.FormatInt(p.GetID(), 10),
					*p.Title,
					*p.URL,
					*p.State,
				)
				// TODO notify
			}
		}

	}

}
