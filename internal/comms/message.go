package comms

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func NewMessage(id, title, url, status string) *Message {
	m := &Message{ID: id,
		Title:  title,
		URL:    url,
		Status: status,
	}
	return m
}

func (m *Message) Format() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("could not format message: %v", err)
	}
	return b
}

// TODO store the data
func (m *Message) Store(redisURL string) error {
	// connect to redis
	// save
}
