package dao

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/stretchr/testify/assert"
)

func TestMessageStatus(t *testing.T) {
	err := d.MessageStatus()
	assert.Nil(t, err)
}
func TestDeviceRegister(t *testing.T) {
	snPrefix := "KKSNAABBCCDDEE%02d%02d"
	keyPrefix := "AABBCCDDEEFFGGHHIIJJKKLLMMNN%02d%02d"

	var s sync.WaitGroup
	for j := int32(0); j < 100; j++ {
		s.Add(1)
		go func(p int32) {
			for i := int32(0); i < 100; i++ {
				sn := fmt.Sprintf(snPrefix, p, i)
				key := fmt.Sprintf(keyPrefix, p, i)
				de := &model.Device{
					Sn:         sn,
					Key:        key,
					Online:     false,
					Platform:   551,
					Serias:     i%2 + 1,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				err := d.DeviceRegister(de)
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
	err := d.MessageReceived(context.TODO(), 1, 14)
	assert.Nil(t, err)
}

func TestDeviceAuthOnline(t *testing.T) {
	var (
		token = model.AuthToken{
			Key: "6aacbf4e43374ad2ac00653de3100a98",
		}
	)
	device, err := d.DeviceAuthOnline(&token)
	println(device)
	println(err)
	assert.Nil(t, err)
}
