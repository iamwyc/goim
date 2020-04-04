package dao

import (
	"errors"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const dbname = "kkgoim"
const deviceCollection = "device"
const offlineMessageCollection = "offline_message"
const messageCollection = "message"

func (d *Dao) NewMessage(message *model.Message) error {
	err := d.getCollection(messageCollection).Insert(message)
	if err != nil {
		return err
	}
	return d.batchInsertDimensionOfflineMessage(message)
}

func (d *Dao) UserRegister(device *model.Device) error {
	return d.getCollection(deviceCollection).Insert(device)
}

func (d *Dao) getCollection(collectionName string) *mgo.Collection {
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

	if m.Online >= 0 {
		dimension["online"] = m.Online == 1
	}
	dCol := d.getCollection(deviceCollection)
	var result []model.Device
	dCol.Find(dimension).Select(bson.M{"id": 1}).All(&result)
	if result == nil || len(result) == 0 {
		return nil
	}
	var messages []interface{}
	for _, r := range result {
		messages = append(messages, model.OfflineMessage{
			Id:       bson.NewObjectId(),
			Seq:      m.Seq,
			DeviceId: r.Id.Hex(),
			Received: 0,
		})
	}

	return d.getCollection(offlineMessageCollection).Insert(messages...)
}

func (d *Dao) MessageReceived(mid string, seq int32) error {
	collection := d.getCollection(offlineMessageCollection)
	_, err := collection.Upsert(bson.M{"deviceId": mid, "seq": seq}, bson.M{"$inc": bson.M{"received": 1}})
	return err
}

func (d *Dao) GetUserOfflineMessage(mid string) (error, []model.OfflineMessageOutVO) {
	var (
		err      error
		seqs     []int32
		messages []model.OfflineMessageOutVO
	)
	omCol := d.getCollection(offlineMessageCollection)
	err = omCol.Find(bson.M{"deviceId": mid, "received": bson.M{"$eq": 0}}).Select(bson.M{"seq": 1}).Distinct("seq", &seqs)
	if err == nil {
		for _, seq := range seqs {
			messages = append(messages, model.OfflineMessageOutVO{
				Seq:     seq,
				Message: uuid.New().String(),
			})
		}
	}
	return err, messages
}
