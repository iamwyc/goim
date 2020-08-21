package dao

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// offlineMessageCollection offline_message collection name
	offlineMessageCollection = "offline_message"
	// messageCollection message collection name
	messageCollection  = "message"
	dbname             = "kkgoim"
	deviceCollection   = "device"
	sequenceCollection = "sequence"
	deviceIDKey        = "device_id"
	messageIDKey       = "message_id"
)

//DeviceAuthOnline get device by token
func (d *Dao) DeviceAuthOnline(token *model.AuthToken) (*model.Device, error) {
	var (
		device model.Device
		change = mgo.Change{
			Update:    bson.M{"$set": bson.M{"online": true}},
			Upsert:    false,
			ReturnNew: true,
		}
	)
	session := d.MongoSession.Copy()
	defer session.Close()

	_, err := d.GetCollection(session, deviceCollection).Find(bson.M{"key": token.Key}).Apply(change, &device)
	return &device, err
}

//DeviceOffline get device by token
func (d *Dao) DeviceOffline(mid int64) error {
	var (
		update = bson.M{"$set": bson.M{"online": false}}
	)
	session := d.MongoSession.Copy()
	defer session.Close()
	return d.GetCollection(session, deviceCollection).Update(bson.M{"_id": mid}, update)
}

// NewMessage insert a new messagepush
func (d *Dao) NewMessage(message *model.Message) (err error) {
	session := d.MongoSession.Copy()
	defer session.Close()
	message.ID = d.idWorker.GetID()
	message.CreateTime = time.Now()
	if err != nil {
		return err
	}
	err = d.GetCollection(session, messageCollection).Insert(message)
	if err != nil {
		return err
	}
	return d.BatchInsertDimensionOfflineMessage(message)
}

// MessageStats receiveed message status
func (d *Dao) MessageStats() (err error) {
	q1 := bson.M{
		"$match": bson.M{
			"received": bson.M{"$gt": 0},
		},
	}
	q2 := bson.M{
		"$group": bson.M{
			"_id":   "$seq",
			"count": bson.M{"$sum": 1},
		},
	}
	q3 := bson.M{
		"$project": bson.M{
			"_id":   1,
			"count": 1,
		},
	}
	var res []model.MessageAggregate
	operations := []bson.M{q1, q2, q3}
	session := d.MongoSession.Copy()
	defer session.Close()
	err = d.GetCollection(session, offlineMessageCollection).Pipe(operations).All(&res)
	if err == nil && len(res) > 0 {
		duration, _ := time.ParseDuration("-72h")
		startTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		startTime = startTime.Add(duration)
		up := bson.M{}
		se := bson.M{"createTime": bson.M{"$gt": startTime}}
		messageCol := d.GetCollection(session, messageCollection)
		for _, v := range res {
			se["_id"] = v.Seq
			up["$set"] = bson.M{"pushCount": v.Count}
			err = messageCol.Update(se, up)
			if err != nil {
				glog.Infof("messageStats update(error %v) : %v", err, v)
			}
		}
	}
	glog.Infof("messageStats(error %v) : %v\n", err, res)
	return
}

// MessageAddSnFile insert a new messagepush
func (d *Dao) MessageAddSnFile(messageID int64, fileList []string) (message *model.Message, err error) {
	var (
		change = mgo.Change{
			Update:    bson.M{"$push": bson.M{"snfile": fileList}},
			Upsert:    false,
			ReturnNew: true,
		}
	)
	session := d.MongoSession.Copy()
	defer session.Close()
	_, err = d.GetCollection(session, messageCollection).Find(bson.M{"_id": messageID}).Apply(change, &message)
	return message, err
}

// DeviceRegister device register
func (d *Dao) DeviceRegister(device *model.Device) (err error) {
	var existDevice model.Device
	session := d.MongoSession.Copy()
	defer session.Close()
	err = d.GetCollection(session, deviceCollection).Find(bson.M{"sn": device.Sn}).One(&existDevice)
	if existDevice.Key != "" {
		device.Key = existDevice.Key
		device.ID = existDevice.ID
		return nil
	}

	device.ID, err = d.getNextSeq(session, deviceIDKey)
	if err != nil {
		return err
	}
	if "" == device.Key {
		device.Key = strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
	}
	device.CreateTime = time.Now()
	device.UpdateTime = time.Now()
	device.Online = false
	return d.GetCollection(session, deviceCollection).Insert(device)
}

// DeviceCount device register
func (d *Dao) DeviceCount() (int, error) {
	session := d.MongoSession.Copy()
	defer session.Close()
	return d.GetCollection(session, deviceCollection).Count()
}

// GetDeviceBySn device register
func (d *Dao) GetDeviceBySn(sn string) (*model.Device, error) {
	var device model.Device
	session := d.MongoSession.Copy()
	defer session.Close()
	err := d.GetCollection(session, deviceCollection).Find(bson.M{"sn": sn}).One(&device)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	return &device, err
}

// GetCollection get mongodb collection by name
func (d *Dao) GetCollection(session *mgo.Session, collectionName string) *mgo.Collection {
	return session.DB(dbname).C(collectionName)
}

// GetMessageByID get message by id
func (d *Dao) GetMessageByID(id int64) (message *model.Message, err error) {
	session := d.MongoSession.Copy()
	defer session.Close()
	err = d.GetCollection(session, messageCollection).Find(bson.M{"_id": id}).One(&message)
	return
}

// GetOfflineMessageByMID get offlinemessage by mid
func (d *Dao) GetOfflineMessageByMID(mid int64) (msgIDList []int64, err error) {
	session := d.MongoSession.Copy()
	defer session.Close()
	omCol := d.GetCollection(session, offlineMessageCollection)
	err = omCol.Find(bson.M{"deviceId": mid, "online": 0, "received": bson.M{"$eq": 0}}).Select(bson.M{"msgId": 1}).Distinct("msgId", &msgIDList)
	return
}

// BatchInsertDimensionOfflineMessage 批量插入
func (d *Dao) BatchInsertDimensionOfflineMessage(m *model.Message) error {

	if m == nil {
		return errors.New("插入维度不能为空")
	}

	timeStr := time.Now().Format("2006-01-02")
	duration, _ := time.ParseDuration("75h")
	startTime, _ := time.Parse("2006-01-02", timeStr)
	expiretTime := startTime.Add(duration)
	var dimension = bson.M{}
	if m.Platform > 0 {
		dimension["platform"] = m.Platform
	}
	if m.Serias > 0 {
		dimension["serias"] = m.Serias
	}
	if m.Sn != nil && len(m.Sn) > 0 {
		dimension["sn"] = bson.M{"$in": m.Sn}
	}
	if m.Online > 0 {
		dimension["online"] = m.Online == 1
	}
	if m.Mids != nil && len(m.Mids) > 0 {
		dimension["_id"] = bson.M{"$in": m.Mids}
	}

	session := d.MongoSession.Copy()
	defer session.Close()
	dCol := d.GetCollection(session, deviceCollection)
	var (
		result   []model.Device
		messages []interface{}
		err      error
		l        = 0
	)
	dCol.Find(dimension).Select(bson.M{"id": 1}).All(&result)
	if result == nil || len(result) == 0 {
		return nil
	}
	for _, r := range result {
		l++
		messages = append(messages, model.OfflineMessage{
			ID:         bson.NewObjectId(),
			Seq:        m.Seq,
			DeviceID:   r.ID,
			Online:     m.Online,
			Received:   0,
			ExpireTime: expiretTime,
		})
		if l%1000 == 0 {
			err = d.GetCollection(session, offlineMessageCollection).Insert(messages...)
			if err != nil {
				glog.Errorf("批量插入离线消息有错误 %v", err)
			}
			messages = messages[0:0]
			l = 0
		}
	}
	if l > 0 {
		err = d.GetCollection(session, offlineMessageCollection).Insert(messages...)
	}
	return err
}

// MessageReceived message received operation
func (d *Dao) MessageReceived(c context.Context, mid int64, msgID int64) error {
	d.MessageSeqAdd(c, msgID)
	session := d.MongoSession.Copy()
	defer session.Close()
	collection := d.GetCollection(session, offlineMessageCollection)
	_, err := collection.UpdateAll(bson.M{"deviceId": mid, "msgId": msgID}, bson.M{"$inc": bson.M{"received": 1}})
	return err
}

type sequence struct {
	NextSeq int32 `bson:"nextSeq"`
}

func (d *Dao) getNextSeq(session *mgo.Session, id string) (int32, error) {
	var (
		seq = sequence{
			NextSeq: int32(1),
		}
		change = mgo.Change{
			Update:    bson.M{"$inc": bson.M{"nextSeq": seq.NextSeq}},
			Upsert:    true,
			ReturnNew: true,
		}
	)
	collection := d.GetCollection(session, sequenceCollection)
	_, err := collection.Find(bson.M{"_id": id}).Apply(change, &seq)
	return seq.NextSeq, err
}
