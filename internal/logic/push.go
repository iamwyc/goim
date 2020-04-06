package logic

import (
	"context"

	"github.com/Terry-Mao/goim/internal/logic/model"

	log "github.com/golang/glog"
)

// PushKeys push a message by keys.
func (l *Logic) PushKeys(c context.Context, arg *model.PushKeyMessage, msg []byte) (err error) {
	message := model.Message{
		Type:      0,
		Online:    arg.Online,
		Operation: arg.Op,
		Content:   msg,
		Sn:        arg.Keys,
	}
	err = l.dao.NewMessage(&message)
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	arg.Seq = message.Seq
	servers, err := l.dao.ServersByKeys(c, arg.Keys)
	if err != nil {
		return
	}
	pushKeys := make(map[string][]string)
	for i, key := range arg.Keys {
		server := servers[i]
		if server != "" && key != "" {
			pushKeys[server] = append(pushKeys[server], key)
		}
	}
	for server := range pushKeys {
		if err = l.dao.PushMsg(c, arg.Op, server, pushKeys[server], arg.Seq, msg); err != nil {
			return
		}
	}
	return
}

//PushMids :push a message by mid.
func (l *Logic) PushMids(c context.Context, arg *model.PushMidsMessage, msg []byte) (err error) {
	message := model.Message{
		Type:      1,
		Online:    arg.Online,
		Operation: arg.Op,
		Content:   msg,
		Mids:      arg.Mids,
	}
	err = l.dao.NewMessage(&message)
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	arg.Seq = message.Seq
	return l.DoPushMids(c, arg, msg)
}

//DoPushMids :do push a message by mid.
func (l *Logic) DoPushMids(c context.Context, arg *model.PushMidsMessage, msg []byte) (err error) {
	keyServers, _, err := l.dao.KeysByMids(c, arg.Mids)
	if err != nil {
		return
	}
	keys := make(map[string][]string)
	for key, server := range keyServers {
		if key == "" || server == "" {
			log.Warningf("push key:%s server:%s is empty", key, server)
			continue
		}
		keys[server] = append(keys[server], key)
	}
	for server, keys := range keys {
		if err = l.dao.PushMsg(c, arg.Op, server, keys, arg.Seq, msg); err != nil {
			return
		}
	}
	return
}

// PushRoom push a message by room.
func (l *Logic) PushRoom(c context.Context, arg *model.PushRoomMessage, msg []byte) (err error) {
	room := model.EncodePlatformRoomKey(arg.Platform)
	message := model.Message{
		Type:      2,
		Online:    arg.Online,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Platform:  arg.Platform,
		Serias:    arg.Serias,
		Content:   msg,
		Room:      room,
	}

	err = l.dao.NewMessage(&message)
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	arg.Seq = message.Seq
	return l.dao.BroadcastRoomMsg(c, arg, room, msg)
}

// PushAll push a message to all.
func (l *Logic) PushAll(c context.Context, arg *model.PushAllMessage, msg []byte) (err error) {
	message := model.Message{
		Type:      3,
		Online:    arg.Online,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Content:   msg,
	}

	err = l.dao.NewMessage(&message)
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	arg.Seq = message.Seq
	return l.dao.BroadcastMsg(c, arg, msg)
}
