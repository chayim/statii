package plugins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/chayim/statii/internal/comms"
	log "github.com/sirupsen/logrus"
)

type JiraConfig struct {
	Username string      `yaml:"username"`
	Token    string      `yaml:"token"`
	Endpoint string      `yaml:"endpoint"`
	Issues   []JiraIssue `yaml:"issues"`
}

type JiraIssue struct {
	Query string `yaml:"query"`
	PluginBase
}

func (j *JiraIssue) Gather(ctx context.Context, user string, pass string, endpoint string, since time.Time) []*comms.Message {
	ts := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	client, err := jira.NewClient(ts.Client(), endpoint)
	if err != nil {
		log.Errorf("error on %s: %v", j.Name, err)
		return nil
	}

	opts := jira.SearchOptions{
		MaxResults: 50,
	}

	issues, resp, err := client.Issue.Search(j.Query, &opts)
	if err != nil {
		log.Errorf("error on %s: %v", j.Name, err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Debugf("no results returned from jira %s", j.Name)
		return nil
	}

	baseURL, _ := url.Parse(endpoint)
	source := fmt.Sprintf("Jira %s", j.Name)
	var messages []*comms.Message
	for _, issue := range issues {
		u := baseURL.JoinPath("browse", issue.ID)
		if !time.Time(issue.Fields.Updated).After(since) {
			continue
		}
		m := comms.NewMessage(
			issue.ID,
			source,
			issue.Fields.Summary,
			u.String(),
			issue.Fields.Status.Name,
		)
		messages = append(messages, m)
	}

	return messages

}
