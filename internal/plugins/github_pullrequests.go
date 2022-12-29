package plugins

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/google/go-github/v48/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

var GITHUB_PR_CONFIG string = "github_pullrequest"

type GitHubPullRequestConfig struct {
	Token        string   `yaml:"token" validate:"required"`
	Repositories []string `yaml:"repositories" validate:"required"`
	States       []string `yaml:"states,omitempty"`
	PluginBase
}

// Gather collects issues, returning items matching the filter, updated
// since the date provided
func (g *GitHubPullRequestConfig) Gather(ctx context.Context, since time.Time) []*comms.Message {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	var messages []*comms.Message

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

		source := fmt.Sprintf("Github [%s]", parts[1])
		for _, p := range pulls {
			if !p.UpdatedAt.After(since) {
				continue
			}
			if len(g.States) != 0 {
				if slices.Contains(g.States, *p.State) {
					m := comms.NewMessage(strconv.FormatInt(*p.ID, 10), source, *p.Title, *p.URL, *p.State)
					messages = append(messages, m)
				}
			} else {
				m := comms.NewMessage(strconv.FormatInt(*p.ID, 10), source, *p.Title, *p.URL, *p.State)
				messages = append(messages, m)
			}
		}

	}
	return messages
}
