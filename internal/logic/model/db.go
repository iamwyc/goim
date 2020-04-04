package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Device struct {
	Id         int32     `bson:"_id"`
	Sn         string    `bson:"sn"`
	Key        string    `bson:"key"`
	Platform   int32     `bson:"platform"`
	Serias     int32     `bson:"serias"`
	Online     bool      `bson:"online"`
	UpdateTime time.Time `bson:"update_time"`
	CreateTime time.Time `bson:"create_time"`
}

type Dimension struct {
	Sn         []string
	PlatformId int32
	SeriasId   int32
	Online     int
}

type OfflineMessage struct {
	Id       bson.ObjectId `bson:"_id"`
	DeviceId int32         `bson:"deviceId"`
	Seq      int32         `bson:"seq"`
	Received int32         `bson:"received"`
}

type OfflineMessageOutVO struct {
	Seq     int32
	Message string
}

type Message struct {
	Id        int32    `bson:"_id"`
	Type      int      `bson:"type"`
	Seq       int32    `bson:"seq"`
	Operation int32    `bson:"opration"`
	Content   []byte   `bson:"content"`
	Sn        []string `bson:"snList"`
	Platform  int32    `bson:"platform"`
	Serias    int32    `bson:"serias"`
	Online    int      `bson:"online"`
	Room      string   `bson:"room"`
	Mids      []int64  `bson:"mids"`
}
