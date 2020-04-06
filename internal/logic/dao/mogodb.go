package dao

import (
	"errors"

	"github.com/Terry-Mao/goim/internal/logic/model"
	log "github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// OfflineMessageCollection offline_message collection name
	OfflineMessageCollection = "offline_message"
	// MessageCollection message collection name
	MessageCollection  = "message"
	dbname             = "kkgoim"
	deviceCollection   = "device"
	sequenceCollection = "sequence"
	deviceIDKey        = "device_id"
	messageIDKey       = "message_id"
)
//GetDevice get device by token
func (d *Dao) GetDevice(token *model.AuthToken) (*model.Device, error) {
	var device model.Device
	err := d.GetCollection(deviceCollection).Find(bson.M{"key":token.Key}).One(&device)
	return &device,err
}

// NewMessage insert a new messagepush
func (d *Dao) NewMessage(message *model.Message) (err error) {
	message.Seq, err = d.getNextSeq(messageIDKey)
	message.ID = message.Seq
	if err != nil {
		return err
	}
	err = d.GetCollection(MessageCollection).Insert(message)
	if err != nil {
		return err
	}
	return d.batchInsertDimensionOfflineMessage(message)
}

// DeviceRegister device register
func (d *Dao) DeviceRegister(device *model.Device) (err error) {
	device.ID, err = d.getNextSeq(deviceIDKey)
	if err != nil {
		return err
	}
	return d.GetCollection(deviceCollection).Insert(device)
}

// GetCollection get mongodb collection by name
func (d *Dao) GetCollection(collectionName string) *mgo.Collection {
	return d.mSession.DB(dbname).C(collectionName)
}

func (d *Dao) batchInsertDimensionOfflineMessage(m *model.Message) error {

	if m == nil {
		return errors.New("插入维度不能为空")
	}
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
	dCol := d.GetCollection(deviceCollection)
	var result []model.Device
	dCol.Find(dimension).Select(bson.M{"id": 1}).All(&result)
	if result == nil || len(result) == 0 {
		return nil
	}
	var messages []interface{}
	for _, r := range result {
		messages = append(messages, model.OfflineMessage{
			ID:       bson.NewObjectId(),
			Seq:      m.Seq,
			DeviceID: r.ID,
			Online:   m.Online,
			Received: 0,
		})
	}

	log.Infof("%d %v", len(messages), messages)
	return d.GetCollection(OfflineMessageCollection).Insert(messages...)
}

// MessageReceived message received operation
func (d *Dao) MessageReceived(mid int64, seq int32) error {
	collection := d.GetCollection(OfflineMessageCollection)
	_, err := collection.Upsert(bson.M{"deviceId": mid, "seq": seq}, bson.M{"$inc": bson.M{"received": 1}})
	return err
}

type sequence struct {
	NextSeq int32 `bson:"nextSeq"`
}

func (d *Dao) getNextSeq(name string) (int32, error) {
	seq := sequence{
		NextSeq: int32(1),
	}
	collection := d.GetCollection(sequenceCollection)
	_, err := collection.Find(bson.M{"_id": name}).Apply(mgo.Change{Update: bson.M{"$inc": bson.M{"nextSeq": seq.NextSeq}}, Upsert: true, ReturnNew: true}, &seq)
	return seq.NextSeq, err
}
