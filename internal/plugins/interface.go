package plugins

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// representing the configuration file
type Config struct {
	Plugins Plugin `yaml:"plugins"`
}

// plugins
type Plugin struct {
	GitHubIssues       []GitHubIssueConfig       `yaml:"github_issues,omitempty"`
	GitHubPullRequests []GitHubPullRequestConfig `yaml:"github_pullrequests,omitempty"`
	JiraIssues         []JiraIssue               `yaml:"jira_issues,omitempty"`
}

type PluginBase struct {
	Timeout int    `yaml:"timeout" validate:"required number"`
	Name    string `yaml:"name" validate:"required"`
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
