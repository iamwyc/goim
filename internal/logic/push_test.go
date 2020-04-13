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
			Op:     int32(900),
			SnList: []string{"test_key"},
			Seq:    900,
		}
	)
	msgID, err := lg.PushSnList(c, &arg, msg)
	println(msgID)
	assert.Nil(t, err)
}

func TestPushMids(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushMidsMessage{
			Op:      int32(900),
			MidList: []int64{12, 13},
			Seq:     900,
		}
	)
	msgID, err := lg.PushMidList(c, &arg, msg)
	println(msgID)
	assert.Nil(t, err)
}

func TestPushRoom(t *testing.T) {
	var (
		c   = context.TODO()
		msg = []byte("hello")
		arg = model.PushRoomMessage{
			Op:       int32(900),
			Platform: 10,
			Serias:   0,
			Seq:      900,
		}
	)
	_, err := lg.PushRoom(c, &arg, msg)
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
	msgID, err := lg.PushAll(c, &arg, msg)
	println(msgID)
	assert.Nil(t, err)
}
