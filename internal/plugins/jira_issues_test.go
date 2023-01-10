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
		query:                "assignee = currentUser()",
		greaterThan:          0,
		lengthNotGreaterThan: 500,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := JiraConfig{Username: username,
				Token:    token,
				Endpoint: endpoint,
			}
			ji := JiraIssue{Query: tc.query}
			messages := ji.Gather(context.TODO(), j.Username,
				j.Token,
				j.Endpoint,
				time.Now().AddDate(0, 0, -100))
			assert.True(t, len(messages) <= tc.lengthNotGreaterThan)
			assert.True(t, len(messages) > tc.greaterThan)
		})
	}
}
