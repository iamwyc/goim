package model

const (
	//DefaultOperation default op
	DefaultOperation = 900
)

//PushMessageIDParam push sn file param
type PushMessageIDParam struct {
	MessageID int32 `form:"messageId" binding:"required"`
}
//PushKeyMessage push message by keys
type PushKeyMessage struct {
	SnList []string `form:"snList"  binding:"required"`
	Online int      `form:"isOnline"`
	Op     int32
	Seq    int32
}

//PushMidsMessage push message by mids
type PushMidsMessage struct {
	MidList []int64 `form:"midList" binding:"required"`
	Online  int     `form:"isOnline"`
	Op      int32
	Seq     int32
}

//PushRoomMessage push message by room
type PushRoomMessage struct {
	Platform int32 `form:"platform" binding:"required"`
	Serias   int32 `form:"serias" binding:"required"`
	Online   int   `form:"isOnline"`
	Op       int32
	Seq      int32
}

//PushAllMessage broadcast message
type PushAllMessage struct {
	Online int `form:"isOnline"`
	Op     int32
	Speed  int32
	Seq    int32
}
