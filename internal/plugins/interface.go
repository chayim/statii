package plugins

import (
	"context"
	"sync"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// representing the configuration file
type Config struct {
	// TODO add defaults
	Plugins         Plugin `yaml:"plugins"`
	RescheduleEvery int    `yaml:"reschedule_seconds" validate:"required number"`
	Database        string `yaml:"database"`
	Size            int64  `yaml:"num_notifications"`
}

// plugins
type Plugin struct {
	GitHubIssues       []GitHubIssueConfig       `yaml:"github_issues,omitempty"`
	GitHubPullRequests []GitHubPullRequestConfig `yaml:"github_pullrequests,omitempty"`
	JiraIssues         []JiraIssue               `yaml:"jira_issues,omitempty"`
}

type PluginBase struct {
	Name string `yaml:"name" validate:"required"`
}

// NewConfigFromBytes returns a Config object, read from the YAML byte stream
func NewConfigFromBytes(b []byte) (*Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	v := validator.New()
	err = v.Struct(cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}

// processPlugins runs through the plugins, storing outputs in the database
// TODO with each process in its own goroutine
func (c *Config) ProcessPlugins() {
	con := comms.NewConnection(c.Database, c.Size)
	ctx := context.TODO()

	since := time.Now().Add(-(time.Second * time.Duration(c.RescheduleEvery)))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var sg sync.WaitGroup
		sg.Add(len(c.Plugins.GitHubIssues))
		for _, g := range c.Plugins.GitHubIssues {
			go func() {
				messages := g.Gather(ctx, since)
				con.SaveMany(ctx, messages)
				sg.Done()
			}()
		}
		wg.Done()
	}()

	// for _, g := range c.Plugins.GitHubPullRequests {
	// 	go func() {
	// 		g.Gather(ctx, since)
	// 		wg.Done()
	// 	}()
	// }

	// for _, g := range c.Plugins.JiraIssues {
	// 	go func() {
	// 		g.Gather(ctx, since)
	// 		wg.Done()
	// 	}()
	// }
}
