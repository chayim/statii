package plugins

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitHubPRs(t *testing.T) {
	var githubToken string = os.Getenv("GITHUB_TOKEN")
	tests := []struct {
		name                 string
		lengthNotGreaterThan int
		greaterThan          int
		repositories         []string
		states               []string
	}{{
		name:                 "no repos",
		lengthNotGreaterThan: 0,
		greaterThan:          -1,
		repositories:         []string{},
		states:               []string{},
	}, {
		name:                 "repositories",
		lengthNotGreaterThan: 500,
		greaterThan:          5,
		repositories:         []string{"redis/redis", "redis/jedis"},
		states:               []string{},
	}, {
		name:                 "repositories and state",
		lengthNotGreaterThan: 50,
		greaterThan:          1,
		repositories:         []string{"redis/redis"},
		states:               []string{"open"},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := githubPRConfig{
				Repositories: tc.repositories,
			}
			messages := g.Gather(context.TODO(), githubToken, time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}

func TestGitHubIssues(t *testing.T) {
	var githubToken string = os.Getenv("GITHUB_TOKEN")
	tests := []struct {
		name                 string
		lengthNotGreaterThan int
		greaterThan          int
		repositories         []string
		states               []string
	}{{
		name:                 "no repos",
		lengthNotGreaterThan: 0,
		greaterThan:          -1,
		repositories:         []string{},
		states:               []string{},
	}, {
		name:                 "repositories",
		lengthNotGreaterThan: 500,
		greaterThan:          5,
		repositories:         []string{"redis/redis", "redis/jedis"},
		states:               []string{},
	}, {
		name:                 "repositories and state",
		lengthNotGreaterThan: 50,
		greaterThan:          1,
		repositories:         []string{"redis/redis"},
		states:               []string{"open"},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := githubIssueConfig{
				Repositories: tc.repositories,
			}
			messages := g.Gather(context.TODO(), githubToken, time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}

func TestGitHubActions(t *testing.T) {
	var githubToken string = os.Getenv("GITHUB_TOKEN")
	tests := []struct {
		name                 string
		lengthNotGreaterThan int
		greaterThan          int
		repository           string
		workflows            []string
		branches             []string
	}{{
		name:                 "no repos",
		lengthNotGreaterThan: 0,
		greaterThan:          -1,
		repository:           "",
		workflows:            []string{},
	}, {
		name:                 "repository",
		lengthNotGreaterThan: 500,
		greaterThan:          5,
		repository:           "redis/redis-py",
		workflows:            []string{"integration.yaml"},
	}, {
		name:                 "repo with branches",
		lengthNotGreaterThan: 100,
		greaterThan:          2,
		repository:           "redis/redis-py",
		workflows:            []string{"integration.yaml"},
		branches:             []string{"master"},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := githubActionConfig{
				Repository:        tc.repository,
				WorkflowFileNames: tc.workflows,
				Branches:          tc.branches,
			}
			messages := g.Gather(context.TODO(), githubToken, time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}

func TestGithubReleases(t *testing.T) {
	var githubToken string = os.Getenv("GITHUB_TOKEN")
	tests := []struct {
		name                 string
		lengthNotGreaterThan int
		greaterThan          int
		repositories         []string
	}{{
		name:                 "no repos!",
		lengthNotGreaterThan: 0,
		greaterThan:          -1,
		repositories:         []string{},
	}, {
		name:                 "one repo",
		lengthNotGreaterThan: 10,
		greaterThan:          0,
		repositories:         []string{"redis/redis-py"},
	}, {
		name:                 "multiple repos",
		lengthNotGreaterThan: 10,
		greaterThan:          2,
		repositories:         []string{"redis/redis-py", "redis/jedis"},
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := githubReleaseConfig{
				Repositories: tc.repositories,
			}
			messages := g.Gather(context.TODO(), githubToken, time.Now().AddDate(0, 0, -30))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}

}
