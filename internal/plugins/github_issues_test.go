package plugins

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitHubIssues(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
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
			g := GitHubIssueConfig{Token: token,
				Repositories: tc.repositories,
			}
			messages := g.Gather(context.TODO(), time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}
