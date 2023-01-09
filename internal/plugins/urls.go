package plugins

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/chayim/statii/internal/comms"
	log "github.com/sirupsen/logrus"
)

type URLConfig struct {
	URLs     []string `yaml:"urls" validate:"required"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	PluginBase
}

func (u *URLConfig) Gather(ctx context.Context, since time.Time) []*comms.Message {
	var messages []*comms.Message
	for idx, url := range u.URLs {
		client := &http.Client{}
		r, _ := http.NewRequest(http.MethodGet, url, nil)
		if u.Username != "" && u.Password != "" {
			r.SetBasicAuth(u.Username, u.Password)
		}
		resp, err := client.Do(r)
		if err != nil {
			log.Errorf("http failure on %s: %v", url, err)
			return nil
		}
		if resp.StatusCode >= 300 {
			log.Errorf("%s returned status code %d", url, resp.StatusCode)
			return nil
		}

		updated := resp.Header.Get("Last-Modified")
		t, err := http.ParseTime(updated)
		if err == nil {
			if !t.After(since) {
				continue
			}
		}
		msg := comms.NewMessage(fmt.Sprintf("%s [%d]", u.Name, idx), url, u.Name, url,
			strconv.FormatInt(int64(resp.StatusCode), 10))
		messages = append(messages, msg)

	}
	return messages
}
