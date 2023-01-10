package plugins

import (
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// representing the configuration file
type Config struct {
	Plugins         Plugin `yaml:"plugins"`
	RescheduleEvery int    `yaml:"reschedule_seconds" validate:"min=60,max=1200"`
	Database        string `yaml:"database"`
	Size            int64  `yaml:"num_notifications" validate:"min=5,max=60"`
}

// plugins
type Plugin struct {
	GitHub     GitHubConfig `yaml:"github,omitempty"`
	URLs       []URLConfig  `yaml:"webpages,omitempty"`
	JiraIssues []JiraIssue  `yaml:"jira_issues,omitempty"`
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
