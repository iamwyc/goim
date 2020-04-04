package logic

import (
	"context"
	"github.com/Terry-Mao/goim/internal/logic/model"

	log "github.com/golang/glog"
)

// PushKeys push a message by keys.
func (l *Logic) PushKeys(c context.Context, arg *model.PushKeyMessage, msg []byte) (err error) {
	err = l.dao.NewMessage(&model.Message{
		Type:      0,
		Online:    -1,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Content:   msg,
		Sn:        arg.Keys,
	})
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	log.Infof("%v", arg)
	servers, err := l.dao.ServersByKeys(c, arg.Keys)
	log.Infof("%v", servers)
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

// PushMids push a message by mid.
func (l *Logic) PushMids(c context.Context, arg *model.PushMidsMessage, msg []byte) (err error) {
	err = l.dao.NewMessage(&model.Message{
		Type:      1,
		Online:    -1,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Content:   msg,
		Mids:      arg.Mids,
	})
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
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
	arg.Room = model.EncodeRoomKey(arg.Type, arg.Room)
	err = l.dao.NewMessage(&model.Message{
		Type:      2,
		Online:    -1,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Content:   msg,
		Room:      arg.Room,
	})
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	return l.dao.BroadcastRoomMsg(c, arg, msg)
}

// PushAll push a message to all.
func (l *Logic) PushAll(c context.Context, arg *model.PushAllMessage, msg []byte) (err error) {
	err = l.dao.NewMessage(&model.Message{
		Type:      3,
		Online:    -1,
		Operation: arg.Op,
		Seq:       arg.Seq,
		Content:   msg,
	})
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	return l.dao.BroadcastMsg(c, arg, msg)
}
