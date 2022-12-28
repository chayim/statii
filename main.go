package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/chayim/statii/internal/plugins"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func gather(cfg string) (*plugins.Config, error) {
	config := plugins.Config{}
	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return &config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &config, err
}

func act(cfg *plugins.Config) {
	ctx := context.TODO()
	for _, p := range cfg.Plugins {
		if p.Type == plugins.GITHUB_ISSUE_CONFIG {
			gh := plugins.NewGitHubIssueConfig(p.Config)
			// TODO you get a thread
			gh.Execute(ctx)
		} else if p.Type == plugins.GITHUB_PR_CONFIG {
			gh := plugins.NewGitHubPullRequestConfig(p.Config)
			// TODO you get a thread
			gh.Execute(ctx)
		}
	}
}

func main() {
	var config string
	var plugins cli.StringSlice

	app := &cli.App{
		Name:  "statii",
		Usage: "do status things",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases:     []string{"c"},
				Name:        "config",
				Usage:       "path to configuration file",
				Destination: &config,
			},
			&cli.StringSliceFlag{
				Aliases:     []string{"p"},
				Name:        "plugin",
				Usage:       "plugins to load",
				Destination: &plugins,
			},
		},
		Action: func(ctx *cli.Context) error {
			if config == "" {
				config = "statii.conf"
			}

			_, err := os.Stat(config)
			if err != nil {
				log.Fatalf("%s does not exist", config)
			}
			cfg := gather(config)
			act(cfg)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
