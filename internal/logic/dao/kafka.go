package dao

import (
	"context"
	"encoding/binary"
	"strconv"

	"github.com/Terry-Mao/goim/internal/logic/model"

	pb "github.com/Terry-Mao/goim/api/logic/grpc"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"gopkg.in/Shopify/sarama.v1"
)

// PushMsg push a message to databus.
func (d *Dao) PushMsg(c context.Context, op int32, server string, keys []string, msgID int64, msg []byte) (err error) {

	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_PUSH,
		Operation: op,
		Server:    server,
		Keys:      keys,
		Msg:       putMessageID(msgID, msg),
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(keys[0]),
		Topic: d.c.Kafka.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(push pushMsgId:%v) error(%v)", msgID, err)
	}
	return
}

// BroadcastRoomMsg push a message to databus.
func (d *Dao) BroadcastRoomMsg(c context.Context, arg *model.PushRoomMessage, room string, msgID int64, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM,
		Operation: arg.Op,
		Room:      room,
		Msg:       putMessageID(msgID, msg),
		Seq:       arg.Seq,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(room),
		Topic: d.c.Kafka.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(broadcast_room pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}

// BroadcastMsg push a message to databus.
func (d *Dao) BroadcastMsg(c context.Context, arg *model.PushAllMessage, msgID int64, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Type:      pb.PushMsg_BROADCAST,
		Operation: arg.Op,
		Speed:     arg.Speed,
		Msg:       putMessageID(msgID, msg),
		Seq:       arg.Seq,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(strconv.FormatInt(int64(arg.Op), 10)),
		Topic: d.c.Kafka.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(broadcast pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}

func putMessageID(msgID int64, content []byte) []byte {

	var msgIDBuf = make([]byte, 8)
	binary.BigEndian.PutUint64(msgIDBuf, uint64(msgID))
	return append(msgIDBuf, content...)
}
