package plugins

import (
	"context"
	"sync"
	"time"

	"github.com/chayim/statii/internal/comms"
	log "github.com/sirupsen/logrus"
)

// processPlugins runs through the plugins, storing outputs in the database
func (c *Config) ProcessPlugins() {
	con := comms.NewConnection(c.Database, c.Size)
	ctx := context.TODO()

	since := time.Now().Add(-(time.Second * time.Duration(c.RescheduleEvery)))

	var wg sync.WaitGroup

	if len(c.Plugins.GitHub.Subscriptions.Issues) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing github issues")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHub.Subscriptions.Issues))
			for _, g := range c.Plugins.GitHub.Subscriptions.Issues {
				go func() {
					messages := g.Gather(ctx, c.Plugins.GitHub.Token, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.GitHub.Subscriptions.PullRequests) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing pull requests")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHub.Subscriptions.PullRequests))
			for _, g := range c.Plugins.GitHub.Subscriptions.PullRequests {
				go func() {
					messages := g.Gather(ctx, c.Plugins.GitHub.Token, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.GitHub.Subscriptions.Actions) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing pull requests")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHub.Subscriptions.Actions))
			for _, g := range c.Plugins.GitHub.Subscriptions.Actions {
				go func() {
					messages := g.Gather(ctx, c.Plugins.GitHub.Token, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.GitHub.Subscriptions.Releases) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing releases")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.GitHub.Subscriptions.Releases))
			for _, g := range c.Plugins.GitHub.Subscriptions.Releases {
				go func() {
					messages := g.Gather(ctx, c.Plugins.GitHub.Token, since)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	if len(c.Plugins.URLs) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing web pages")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.URLs))
			for _, g := range c.Plugins.URLs {
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

	if len(c.Plugins.Jira.Issues) > 0 {
		wg.Add(1)
		go func() {
			log.Debug("processing jira")
			var sg sync.WaitGroup
			sg.Add(len(c.Plugins.Jira.Issues))
			for _, g := range c.Plugins.Jira.Issues {
				go func() {
					messages := g.Gather(ctx,
						c.Plugins.Jira.Username,
						c.Plugins.Jira.Token,
						c.Plugins.Jira.Endpoint,
						since,
					)
					con.SaveMany(ctx, messages)
					sg.Done()
				}()
			}
			sg.Wait()
			wg.Done()
		}()
	}

	wg.Wait()

}
