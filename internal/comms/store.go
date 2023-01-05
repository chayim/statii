package comms

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	DB         *redis.Client
	StreamName string
	Size       int64
}

func NewConnection(conn string, size int64) *Store {
	dbcon := redis.NewClient(&redis.Options{
		Addr: conn,
		DB:   0,
	})
	st := Store{StreamName: "stattistream", DB: dbcon, Size: size}
	return &st
}

func (s *Store) Clear(ctx context.Context) error {
	_, err := s.DB.Del(ctx, s.StreamName).Result()
	return err
}

func (s *Store) SaveMany(ctx context.Context, messages []*Message) error {

	if len(messages) == 0 {
		return nil
	}

	pipe := s.DB.Pipeline()
	for _, msg := range messages {
		bytes, _ := json.Marshal(msg)
		pipe.LPush(ctx, s.StreamName, bytes)
	}
	pipe.LTrim(ctx, s.StreamName, 0, s.Size-1)
	log.Debugf("saving to redis at %s", s.DB)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *Store) ReadAll(ctx context.Context) ([]*Message, error) {
	log.Debug("reading from redis")
	res, err := s.DB.LPopCount(ctx, s.StreamName, int(s.Size)).Result()
	if err != nil {
		return nil, err
	}

	var messages []*Message
	for _, r := range res {
		var m *Message
		err := json.Unmarshal([]byte(r), &m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
