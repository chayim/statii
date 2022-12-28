package plugins

// representing the configuration file
type Config struct {
	Timeout int      `yaml:"timeout"`
	Plugins []Plugin `yaml:"plugins"`
	Redis   string   `yaml:"redis"`
}

// plugins
type Plugin struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"config"`
	// GitHubIssues       []GitHubIssueConfig `yaml:"github_issues"`
	// GitHubPullRequests []GitHubPullRequest `yaml:"github_pullrequests"`
}
