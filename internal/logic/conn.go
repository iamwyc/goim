package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/dao"
	"gopkg.in/mgo.v2/bson"

	"github.com/Terry-Mao/goim/api/comet/grpc"
	"github.com/Terry-Mao/goim/internal/logic/model"
	log "github.com/golang/glog"
	"github.com/google/uuid"
)

type AuthToken struct {
	Mid      int32   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

// Connect connected a conn.
func (l *Logic) Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error) {
	var params AuthToken
	if err = json.Unmarshal(token, &params); err != nil {
		log.Errorf("json.Unmarshal(%s) error(%v)", token, err)
		return
	}
	mid = int64(params.Mid)
	roomID = params.RoomID + "@test://123"
	accepts = params.Accepts
	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	if key = params.Key; key == "" {
		key = uuid.New().String()
	}
	if err = l.dao.AddMapping(c, mid, key, server); err != nil {
		log.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
	}
	log.Infof("conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

// Disconnect disconnect a conn.
func (l *Logic) Disconnect(c context.Context, mid int64, key, server string) (has bool, err error) {
	if has, err = l.dao.DelMapping(c, mid, key, server); err != nil {
		log.Errorf("l.dao.DelMapping(%d,%s) error(%v)", mid, key, server)
		return
	}
	log.Infof("conn disconnected key:%s server:%s mid:%d", key, server, mid)
	return
}

// Heartbeat heartbeat a conn.
func (l *Logic) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	has, err := l.dao.ExpireMapping(c, mid, key)
	if err != nil {
		log.Errorf("l.dao.ExpireMapping(%d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	if !has {
		if err = l.dao.AddMapping(c, mid, key, server); err != nil {
			log.Errorf("l.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err)
			return
		}
	}
	log.Infof("conn heartbeat key:%s server:%s mid:%d", key, server, mid)
	return
}

// RenewOnline renew a server online.
func (l *Logic) RenewOnline(c context.Context, server string, roomCount map[string]int32) (map[string]int32, error) {
	online := &model.Online{
		Server:    server,
		RoomCount: roomCount,
		Updated:   time.Now().Unix(),
	}
	if err := l.dao.AddServerOnline(context.Background(), server, online); err != nil {
		return nil, err
	}
	return l.roomCount, nil
}

// Receive receive a message.
func (l *Logic) Receive(c context.Context, mid int64, proto *grpc.Proto) (err error) {
	log.Infof("receive mid:%d message:%+v", mid, proto)
	if proto.Op == grpc.OpGetOfflineMessage {
		go l.GetUserOfflineMessage(mid)
		proto.Op = grpc.OpGetOfflineMessageReply
	} else if proto.Op == grpc.OpBusinessMessageAck {
		l.dao.MessageReceived(mid, proto.Seq)
	}
	return
}

// GetUserOfflineMessage get user offline message Operate
func (l *Logic) GetUserOfflineMessage(mid int64) error {
	var (
		err  error
		seqs []int32
	)
	ctx := context.TODO()
	omCol := l.dao.GetCollection(dao.OfflineMessageCollection)
	err = omCol.Find(bson.M{"deviceId": mid, "received": bson.M{"$eq": 0}}).Select(bson.M{"seq": 1}).Distinct("seq", &seqs)

	if err == nil {
		mids := []int64{mid}
		for _, seq := range seqs {
			var message model.Message
			err := l.dao.GetCollection(dao.MessageCollection).Find(bson.M{"_id": seq}).One(&message)
			if err == nil {
				pm := model.PushMidsMessage{
					Op:   message.Operation,
					Seq:  message.Seq,
					Mids: mids,
				}
				l.DoPushMids(ctx, &pm, message.Content)
			} else {
				log.Errorf("查询消息出错%v", err)
			}
		}
	} else {
		log.Errorf("查询消息出错%v", err)
	}
	return err
}
