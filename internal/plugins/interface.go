package plugins

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/chayim/statii/internal/comms"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// representing the configuration file
type Config struct {
	// TODO add defaults
	Plugins         Plugin `yaml:"plugins"`
	RescheduleEvery int    `yaml:"reschedule_seconds" validate:"min=60,max=1200"`
	Database        string `yaml:"database"`
	Size            int64  `yaml:"num_notifications" validate:"min=5,max=60"`
}

// plugins
type Plugin struct {
	GitHubIssues       []GitHubIssueConfig       `yaml:"github_issues,omitempty"`
	GitHubPullRequests []GitHubPullRequestConfig `yaml:"github_pullrequests,omitempty"`
	GitHubActions      []GitHubActionConfig      `yaml:"github_actions,omitempty"`
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
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			log.Error(e)
		}
		os.Exit(3)
	}
	return &cfg, err
}

// processPlugins runs through the plugins, storing outputs in the database
func (c *Config) ProcessPlugins() {
	con := comms.NewConnection(c.Database, c.Size)
	ctx := context.TODO()

	since := time.Now().Add(-(time.Second * time.Duration(c.RescheduleEvery)))

	var wg sync.WaitGroup

	if len(c.Plugins.GitHubIssues) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing github issues")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHubIssues))
			for _, g := range c.Plugins.GitHubIssues {
				go func() {
					messages := g.Gather(ctx, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.GitHubPullRequests) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing pull requests")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHubPullRequests))
			for _, g := range c.Plugins.GitHubPullRequests {
				go func() {
					messages := g.Gather(ctx, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.GitHubActions) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing pull requests")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHubActions))
			for _, g := range c.Plugins.GitHubActions {
				go func() {
					messages := g.Gather(ctx, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	// for _, g := range c.Plugins.JiraIssues {
	// 	go func() {
	// 		g.Gather(ctx, since)
	// 		wg.Done()
	// 	}()
	// }

	wg.Wait()

}
