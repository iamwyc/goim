package dao

import (
	"fmt"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"testing"
	"time"
)

func TestUserRegister(t *testing.T) {
	for i := int32(10); i < 100; i++ {
		de := &model.Device{
			Id:         bson.NewObjectId(),
			Sn:         "00:11:22:33:44:" + strconv.Itoa(int(i)),
			Key:        uuid.New().String(),
			Online:     false,
			Platform:   i % 10,
			Serias:     i % 4,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		d.UserRegister(de)
	}
}

func TestNewMessage(t *testing.T) {
	d.NewMessage(&model.Message{
		Seq:       100,
		Operation: 1000,
		Content:   []byte("TestNewMessage"),
		Serias:    0,
		Platform:  2,
		Sn:        nil,
		Online:    0,
	})
}
func TestMessageReceived(t *testing.T) {
	d.MessageReceived("5e8847d7f4f2f43b808da48e", 19)
}

func TestGetUserOfflineMessage(t *testing.T) {
	err, messages := d.GetUserOfflineMessage("5e88976bf4f2f44d50f1a897")
	if err == nil {
		fmt.Printf("%v", messages)
	} else {
		panic(err)
	}
}
