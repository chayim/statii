package comms

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var connstr string = "localhost:9876"
var size int64 = 25

func messages(num int) []*Message {
	messages := []*Message{}
	for i := 0; i < num; i++ {
		m := Message{
			ID:     fmt.Sprintf("myid%d", i),
			URL:    "http://foo.bar",
			Title:  "xxx",
			Status: "green",
			Source: "something",
		}
		messages = append(messages, &m)
	}

	return messages
}

func TestConnection(t *testing.T) {
	con := NewConnection(connstr, size)
	ctx := context.Background()

	// start empty
	con.DB.FlushAll(ctx).Result()
	set, err := con.DB.Set(ctx, "key", "value", 10*time.Second).Result()
	assert.Equal(t, "OK", set)
	assert.True(t, err == nil)
	_, err = con.DB.FlushDB(ctx).Result()
	assert.True(t, err == nil)
}

func TestSaveMany(t *testing.T) {
	tests := []struct {
		name   string
		length int
		max    int64
	}{{
		name:   "less than max",
		length: 10,
		max:    10,
	}, {
		name:   "more than max",
		length: 99,
		max:    25,
	}, {
		name:   "equal to max",
		length: 25,
		max:    25,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			connstr = "localhost:6379"
			con := NewConnection(connstr, size)
			ctx := context.Background()

			// start empty
			con.DB.FlushAll(ctx).Result()

			err := con.SaveMany(ctx, messages(tc.length))
			assert.True(t, err == nil)

			res, _ := con.DB.LLen(ctx, con.StreamName).Result()
			assert.Equal(t, tc.max, res)

			con.DB.FlushAll(ctx).Result()
			r, _ := con.DB.Keys(ctx, "*").Result()
			assert.Equal(t, []string{}, r)
		})
	}

}

func TestReadAll(t *testing.T) {

	con := NewConnection(connstr, size)

	ctx := context.Background()
	// start empty
	con.DB.FlushAll(ctx).Result()

	num := 5
	con.SaveMany(ctx, messages(num))
	messages, err := con.ReadAll(ctx)
	assert.True(t, err == nil)

	assert.Len(t, messages, num)
}
