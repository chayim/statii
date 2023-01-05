package comms

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func NewMessage(id, source, title, url, status string) *Message {
	m := &Message{ID: id,
		Source: source,
		Title:  title,
		URL:    url,
		Status: status,
	}
	return m
}

func (m *Message) Dump() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("could not format message: %v", err)
	}
	return b
}
