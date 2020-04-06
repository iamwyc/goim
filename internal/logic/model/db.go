package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//Device device model
type Device struct {
	ID         int32     `bson:"_id"`
	Sn         string    `bson:"sn"`
	Key        string    `bson:"key"`
	Platform   int32     `bson:"platform"`
	Serias     int32     `bson:"serias"`
	Online     bool      `bson:"online"`
	UpdateTime time.Time `bson:"update_time"`
	CreateTime time.Time `bson:"create_time"`
}

//Dimension dimension model
type Dimension struct {
	Sn         []string
	PlatformID int32
	SeriasID   int32
	Online     int
}

//OfflineMessage offline message model
type OfflineMessage struct {
	ID       bson.ObjectId `bson:"_id"`
	DeviceID int32         `bson:"deviceId"`
	Seq      int32         `bson:"seq"`
	Online   int           `bson:"online"`
	Received int32         `bson:"received"`
}

//Message message model
type Message struct {
	ID        int32    `bson:"_id"`
	Type      int      `bson:"type"`
	Seq       int32    `bson:"seq"`
	Operation int32    `bson:"opration"`
	Content   []byte   `bson:"content"`
	Sn        []string `bson:"snList"`
	Platform  int32    `bson:"platform"`
	Serias    int32    `bson:"serias"`
	//0:不限 1:在线消息
	Online int     `bson:"online"`
	Room   string  `bson:"room"`
	Mids   []int64 `bson:"mids"`
}
