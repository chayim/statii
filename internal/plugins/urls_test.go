package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLFetch(t *testing.T) {
	tests := []struct {
		name   string
		page   []string
		length int
	}{{
		name:   "google",
		page:   []string{"https://google.ca"},
		length: 1,
	}, {
		name:   "news",
		page:   []string{"https://news.google.com", "https://news.com"},
		length: 2,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := URLConfig{URLs: tc.page}
			messages := u.Gather(context.TODO(), time.Now().AddDate(0, 0, -15))
			assert.True(t, len(messages) == tc.length)
		})
	}
}
