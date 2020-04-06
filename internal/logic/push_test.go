package logic

import (
	"context"
	"testing"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/stretchr/testify/assert"
)

func TestPushKeys(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushKeyMessage{
			Op:   int32(900),
			Keys: []string{"test_key"},
			Seq:  900,
		}
	)
	err := lg.PushKeys(c, &arg, msg)
	assert.Nil(t, err)
}

func TestPushMids(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushMidsMessage{
			Op:   int32(900),
			Mids: []int64{12, 13},
			Seq:  900,
		}
	)
	err := lg.PushMids(c, &arg, msg)
	assert.Nil(t, err)
}

func TestPushRoom(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushRoomMessage{
			Op:   int32(900),
			Room: "123",
			Type: "test",
			Seq:  900,
		}
	)
	err := lg.PushRoom(c, &arg, msg)
	assert.Nil(t, err)
}

func TestPushAll(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushAllMessage{
			Op:    int32(900),
			Seq:   900,
			Speed: 5,
		}
	)
	err := lg.PushAll(c, &arg, msg)
	assert.Nil(t, err)
}
