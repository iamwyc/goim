package dao

import (
	"fmt"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestUserRegister(t *testing.T) {
	snPrefix := "00:11:22:33:%d:%d"
	var s sync.WaitGroup
	for j := int32(10); j < 100; j++ {
		s.Add(1)
		go func(p int32) {
			for i := int32(10); i < 100; i++ {
				sn := fmt.Sprintf(snPrefix, p, i)
				de := &model.Device{
					Sn:         sn,
					Key:        uuid.New().String(),
					Online:     false,
					Platform:   i % 10,
					Serias:     i % 4,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				err := d.UserRegister(de)
				if err != nil {
					panic(err)
				}
			}
			s.Done()
		}(j)
	}
	s.Wait()
}

func TestNewMessage(t *testing.T) {
	err := d.NewMessage(&model.Message{
		Seq:       100,
		Operation: 1000,
		Content:   []byte("TestNewMessage"),
		Serias:    0,
		Platform:  2,
		Sn:        nil,
		Online:    0,
	})
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)

}
func TestMessageReceived(t *testing.T) {
	d.MessageReceived(1, 19)
}

func TestGetUserOfflineMessage(t *testing.T) {
	err, ms := d.GetUserOfflineMessage(1)
	assert.Nil(t, err)
	fmt.Printf("%v", ms)
}
