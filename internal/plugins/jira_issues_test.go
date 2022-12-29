package plugins

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJira(t *testing.T) {

	token := os.Getenv("JIRA_TOKEN")
	username := os.Getenv("JIRA_USERNAME")
	endpoint := os.Getenv("JIRA_ENDPOINT")

	tests := []struct {
		name                 string
		query                string
		lengthNotGreaterThan int
		greaterThan          int
	}{{
		name:                 "no results query",
		query:                "text~\"fooshoomooogoo\"",
		greaterThan:          -1,
		lengthNotGreaterThan: 0,
	}, {
		name:                 "invalid query",
		query:                "shooboogamoo",
		greaterThan:          -1,
		lengthNotGreaterThan: 0,
	}, {
		name:                 "basic query",
		query:                "text~well",
		greaterThan:          0,
		lengthNotGreaterThan: 500,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := JiraIssue{Username: username,
				Token:    token,
				Endpoint: endpoint,
				Query:    tc.query,
			}
			messages := j.Gather(context.TODO(), time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}
