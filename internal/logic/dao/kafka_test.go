package dao

import (
	"context"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaoPushMsg(t *testing.T) {
	var (
		c      = context.Background()
		op     = int32(100)
		server = "test"
		msg    = []byte("msg")
		keys   = []string{"key"}
	)
	err := d.PushMsg(c, op, server, keys, 1, msg)
	assert.Nil(t, err)
}

func TestDaoBroadcastRoomMsg(t *testing.T) {
	var (
		c   = context.Background()
		msg = []byte("msg")
		arg = model.PushRoomMessage{
			Op:   int32(100),
			Room: "123",
			Seq:  1,
			Type: "test",
		}
	)

	err := d.BroadcastRoomMsg(c, &arg, msg)
	assert.Nil(t, err)
}

func TestDaoBroadcastMsg(t *testing.T) {
	var (
		c   = context.Background()
		msg = []byte("")
		arg = model.PushAllMessage{
			Op:    int32(100),
			Seq:   1,
			Speed: 5,
		}
	)
	err := d.BroadcastMsg(c, &arg, msg)
	assert.Nil(t, err)
}
