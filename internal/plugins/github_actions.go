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
	"golang.org/x/oauth2"
)

type GitHubActionConfig struct {
	Token             string   `yaml:"token" validate:"required"`
	Repository        string   `yaml:"repository" validate:"required"`
	WorkflowFileNames []string `yaml:"workflow_files" validate:"required"`
	Branches          []string `yaml:"branches" validate:"required"`
	PluginBase
}

// Gather github actions occuring since the specified time
func (g *GitHubActionConfig) Gather(ctx context.Context, since time.Time) []*comms.Message {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	var messages []*comms.Message

	parts := strings.Split(g.Repository, "/")
	if len(parts) != 2 {
		log.Debugf("%s is an invalid repository, skipping.", g.Repository)
		return nil
	}

	log.WithFields(log.Fields{"config": g.Name, "repo": g.Repository}).Debug("github actions")
	for _, f := range g.WorkflowFileNames {
		for _, branch := range g.Branches {
			opts := github.ListWorkflowRunsOptions{Branch: branch}
			runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, parts[0], parts[1], f, &opts)
			if err != nil {
				log.Errorf("error on %s: %v", g.Repository, err)
				return nil
			}

			source := fmt.Sprintf("Github Workflow [%s:%s]", parts, f)
			for _, r := range runs.WorkflowRuns {

				// out of date spread
				if !r.UpdatedAt.After(since) {
					break
				}
				m := comms.NewMessage(strconv.FormatInt(*r.ID, 10), source, *r.Name, *r.URL, *r.Status)
				messages = append(messages, m)
			}
		}
	}
	return messages
}
