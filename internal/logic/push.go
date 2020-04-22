package logic

import (
	"bufio"
	"context"
	"errors"
	"os"

	"github.com/Terry-Mao/goim/internal/logic/model"

	"github.com/golang/glog"
	log "github.com/golang/glog"
)

// PushSnList push a message by keys.
func (l *Logic) PushSnList(c context.Context, arg *model.PushKeyMessage, msg []byte) (msgID int32, err error) {
	var (
		message *model.Message
	)
	if arg.MessageID > 0 {
		message, err = l.dao.GetMessageByID(arg.MessageID)
		if err == nil {
			message.Sn = arg.SnList
			err = l.dao.BatchInsertDimensionOfflineMessage(message)
		}
	} else {
		message = &model.Message{
			Type:      0,
			Online:    arg.Online,
			Operation: arg.Op,
			Content:   msg,
			Sn:        arg.SnList,
		}
		err = l.dao.NewMessage(message)
	}
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	msgID = message.ID
	arg.Seq = message.Seq
	err = l.doPushSnList(c, message, arg.SnList)
	return
}
func (l *Logic) doPushSnList(c context.Context, message *model.Message, snList []string) (err error) {
	servers, err := l.dao.ServersByKeys(c, snList)
	if err != nil {
		return
	}
	pushSnList := make(map[string][]string)
	for i, key := range snList {
		server := servers[i]
		if server != "" && key != "" {
			pushSnList[server] = append(pushSnList[server], key)
		}
	}
	for server := range pushSnList {
		if err = l.dao.PushMsg(c, message.Operation, server, pushSnList[server], message.Seq, message.Content); err != nil {
			return
		}
	}
	return
}

//PushMidList :push a message by mid.
func (l *Logic) PushMidList(c context.Context, arg *model.PushMidsMessage, msg []byte) (msgID int32, err error) {
	message := model.Message{
		Type:      1,
		Online:    arg.Online,
		Operation: arg.Op,
		Content:   msg,
		Mids:      arg.MidList,
	}
	err = l.dao.NewMessage(&message)
	if err != nil {
		log.Errorf("插入数据库错误:%v", err)
		return
	}
	arg.Seq = message.Seq
	err = l.DoPushMids(c, arg, msg)
	if err == nil {
		msgID = message.ID
	}
	return
}

//DoPushMids :do push a message by mid.
func (l *Logic) DoPushMids(c context.Context, arg *model.PushMidsMessage, msg []byte) (err error) {
	keyServers, _, err := l.dao.KeysByMids(c, arg.MidList)
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

var (
	// ErrRoomNull 异常:平台系列不能为空
	ErrRoomNull = errors.New("必须指定平台或者系列")
)

// PushRoom push a message by room.
func (l *Logic) PushRoom(c context.Context, arg *model.PushRoomMessage, msg []byte) (id int32, err error) {
	room := model.DecodePlatformAndSeriasRoomKey(arg.Platform, arg.Serias)
	if room == "" {
		return 0, ErrRoomNull
	}
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
	id = message.ID
	l.dao.BroadcastRoomMsg(c, arg, room, msg)
	return
}

// PushAll push a message to all.
func (l *Logic) PushAll(c context.Context, arg *model.PushAllMessage, msg []byte) (msgID int32, err error) {
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
	msgID = message.ID
	arg.Seq = message.Seq
	err = l.dao.BroadcastMsg(c, arg, msg)
	return
}

//PushSnFils push by sn file
func (l *Logic) PushSnFils(messageID int32, fileList []string) {
	glog.Infof("fileList :%v", fileList)
	message, err := l.dao.MessageAddSnFile(messageID, fileList)
	if len(fileList) == 0 || err != nil {
		glog.Errorf("存储文件到mongo失败 %d,%v", messageID, fileList, err)
		return
	}

	for _, file := range fileList {
		err := l.pushFile(message, file)
		if err != nil {
			glog.Errorf("串号文件推送失败:%v %v", file, err)
		}
	}
}

// GetFileDir get save path
func (l *Logic) GetFileDir() string {
	return l.c.MessagePush.Dir
}

func (l *Logic) pushFile(message *model.Message, snfilepath string) (err error) {
	file, err := os.Open(snfilepath)
	if err != nil {
		return
	}
	defer file.Close()
	var (
		snList              []string
		ctx                 = context.TODO()
		num                 = 0
		scanner             = bufio.NewScanner(file)
		offlineMessageParam = model.Message{
			ID:        message.ID,
			Seq:       message.Seq,
			Online:    message.Online,
			Operation: message.Operation,
			Content:   message.Content,
		}
	)

	for scanner.Scan() {
		snList = append(snList, scanner.Text())
		num++
		if 0 == num%l.c.MessagePush.BatchPushCount {
			offlineMessageParam.Sn = snList
			l.pushSnFileList(ctx, &offlineMessageParam)
			snList = snList[0:0]
			num = 0
		}
	}
	if len(snList) != 0 {
		offlineMessageParam.Sn = snList
		l.pushSnFileList(ctx, &offlineMessageParam)
	}

	if err = scanner.Err(); err != nil {
		return
	}
	return
}

// MessageRedisStats redis mesage stats
func (l *Logic) MessageRedisStats(ctx context.Context, messageID int32) (int64, error) {
	return l.dao.MessageCountStats(ctx, messageID)
}
func (l *Logic) pushSnFileList(ctx context.Context, message *model.Message) {
	err := l.dao.BatchInsertDimensionOfflineMessage(message)
	if err != nil {
		glog.Errorf("批量插入sn异常:%v %v", message.Sn, err)
	}
	err = l.doPushSnList(ctx, message, message.Sn)
	if err != nil {
		glog.Errorf("批量推送异常:%v %v", message.Sn, err)
	}
}
