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

type GitHubConfig struct {
	Token         string             `yaml:"token" validate:"required"`
	Subscriptions GitHubSubscription `yaml:"subscriptions" validate:"required"`
}

type GitHubSubscription struct {
	Issues       []githubIssueConfig   `yaml:"issues,omitempty"`
	PullRequests []githubPRConfig      `yaml:"pullrequests,omitempty"`
	Actions      []githubActionConfig  `yaml:"actions,omitempty"`
	Releases     []githubReleaseConfig `yaml:"releases,omitempty"`
}

// github issues
type githubIssueConfig struct {
	Repositories []string `yaml:"repositories" validate:"required"`
	States       []string `yaml:"states"`
	Assignee     string   `yaml:"assignee"`
	PluginBase
}

// pull requests
type githubPRConfig struct {
	Repositories []string `yaml:"repositories" validate:"required"`
	States       []string `yaml:"states,omitempty"`
	PluginBase
}

// actions
type githubActionConfig struct {
	Repository        string   `yaml:"repository" validate:"required"`
	WorkflowFileNames []string `yaml:"workflow_files" validate:"required"`
	Branches          []string `yaml:"branches" validate:"required"`
	PluginBase
}

type githubReleaseConfig struct {
	Repositories []string `yaml:"repositories" validate:"required"`
	PluginBase
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

// Gather collects issues, returning items matching the filter, updated
// since the date provided
func (p *githubPRConfig) Gather(ctx context.Context, token string, since time.Time) []*comms.Message {

	var messages []*comms.Message
	client := newGitHubClient(ctx, token)

	for _, repo := range p.Repositories {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			log.Debugf("%s is an invalid repository, skipping.", repo)
			continue
		}

		opts := github.PullRequestListOptions{State: "all", Sort: "updated", Direction: "desc"}
		pulls, _, err := client.PullRequests.List(ctx, parts[0], parts[1], &opts)
		if err != nil {
			log.Errorf("error on %s: %v", repo, err)
			continue
		}

		source := fmt.Sprintf("Github PR [%s]", parts[1])
		for _, pull := range pulls {
			if !pull.UpdatedAt.After(since) {
				break
			}
			if len(p.States) != 0 {
				if slices.Contains(p.States, *pull.State) {
					m := comms.NewMessage(strconv.FormatInt(*pull.ID, 10), source, *pull.Title, *pull.URL, *pull.State)
					messages = append(messages, m)
				}
			} else {
				m := comms.NewMessage(strconv.FormatInt(*pull.ID, 10), source, *pull.Title, *pull.URL, *pull.State)
				messages = append(messages, m)
			}
		}

	}
	return messages
}

// Gather collects issues, returning items matching the filter, updated
// since the date provided
func (i *githubIssueConfig) Gather(ctx context.Context, token string, since time.Time) []*comms.Message {

	var messages []*comms.Message
	client := newGitHubClient(ctx, token)

	for _, repo := range i.Repositories {
		log.WithFields(log.Fields{"config": i.Name, "repo": repo}).Debug("github issues")
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			log.Debugf("%s is an invalid repository, skipping.", repo)
			continue
		}
		opts := github.IssueListByRepoOptions{State: "all", Sort: "updated", Direction: "desc"}
		if i.Assignee != "" {
			opts.Assignee = i.Assignee
		}
		issues, _, err := client.Issues.ListByRepo(ctx, parts[0], parts[1], &opts)
		if err != nil {
			log.Errorf("error on %s: %v", repo, err)
			continue
		}
		source := fmt.Sprintf("Github Issue [%s]", parts[1])
		for _, p := range issues {
			if !p.UpdatedAt.After(since) {
				break
			}
			if len(i.States) != 0 {
				if slices.Contains(i.States, *p.State) {
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

// Gather github actions occuring since the specified time
func (a *githubActionConfig) Gather(ctx context.Context, token string, since time.Time) []*comms.Message {

	var messages []*comms.Message
	client := newGitHubClient(ctx, token)

	parts := strings.Split(a.Repository, "/")
	if len(parts) != 2 {
		log.Debugf("%s is an invalid repository, skipping.", a.Repository)
		return nil
	}

	log.WithFields(log.Fields{"config": a.Name, "repo": a.Repository}).Debug("github actions")
	for _, f := range a.WorkflowFileNames {
		if len(a.Branches) > 0 {
			for _, branch := range a.Branches {
				opts := github.ListWorkflowRunsOptions{Branch: branch}
				runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, parts[0], parts[1], f, &opts)
				if err != nil {
					log.Errorf("error on %s: %v", a.Repository, err)
					return nil
				}

				source := fmt.Sprintf("Github Workflow [%s:%s %s]", parts, f, branch)
				for _, r := range runs.WorkflowRuns {

					// out of date spread
					if !r.UpdatedAt.After(since) {
						break
					}
					m := comms.NewMessage(strconv.FormatInt(*r.ID, 10), source, *r.Name, *r.URL, *r.Status)
					messages = append(messages, m)
				}
			}
		} else {
			runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, parts[0], parts[1], f, nil)
			if err != nil {
				log.Errorf("error on %s: %v", a.Repository, err)
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

// releases
func (r *githubReleaseConfig) Gather(ctx context.Context, token string, since time.Time) []*comms.Message {
	var messages []*comms.Message
	client := newGitHubClient(ctx, token)
	for _, repo := range r.Repositories {
		log.WithFields(log.Fields{"config": r.Name, "repo": repo}).Debug("github issues")
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			log.Debugf("%s is an invalid repository, skipping.", repo)
			continue
		}

		source := fmt.Sprintf("Github Release [%s]", parts)

		releases, _, err := client.Repositories.ListReleases(ctx, parts[0], parts[1], nil)
		if err != nil {
			log.Errorf("error on %s: %v", repo, err)
			continue
		}
		for _, release := range releases {
			if !release.CreatedAt.After(since) {
				break
			}

			m := comms.NewMessage(*release.Name,
				source, *release.Name, *release.URL, "")
			messages = append(messages, m)
		}

	}

	return messages

}
