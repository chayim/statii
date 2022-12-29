package plugins

// representing the configuration file
type Config struct {
	Timeout int    `yaml:"timeout"`
	Plugins Plugin `yaml:"plugins"`
	Redis   string `yaml:"redis_url"`
}

// plugins
type Plugin struct {
	GitHubIssues       []GitHubIssueConfig       `yaml:"github_issues,omitempty"`
	GitHubPullRequests []GitHubPullRequestConfig `yaml:"github_pullrequests,omitempty"`
	JiraIssues         []JiraIssue               `yaml:"jira_issues,omitempty"`
}

type PluginBase struct {
	Name string `yaml:"name"`
}
