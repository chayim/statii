package plugins

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitHubActions(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	tests := []struct {
		name                 string
		lengthNotGreaterThan int
		greaterThan          int
		repository           string
		workflows            []string
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
		workflows:            []string{"integration.yml"},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := GitHubActionConfig{Token: token,
				Repository:        tc.repository,
				WorkflowFileNames: tc.workflows,
			}
			messages := g.Gather(context.TODO(), time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}
