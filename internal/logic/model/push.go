package model

//PushKeyMessage push message by keys
type PushKeyMessage struct {
	Op     int32    `form:"operation"`
	Keys   []string `form:"keys"`
	Online int      `form:"isOnline"`
	Seq    int32
}

//PushMidsMessage push message by mids
type PushMidsMessage struct {
	Op     int32   `form:"operation"`
	Mids   []int64 `form:"mids"`
	Online int     `form:"isOnline"`
	Seq    int32
}

//PushRoomMessage push message by room
type PushRoomMessage struct {
	Op       int32  `form:"operation" binding:"required"`
	Platform int32 `form:"platform" binding:"required"`
	Serias   int32 `form:"serias" binding:"required"`
	Online   int    `form:"isOnline"`
	Seq      int32
}

//PushAllMessage broadcast message
type PushAllMessage struct {
	Op     int32 `form:"operation" binding:"required"`
	Speed  int32 `form:"speed"`
	Online int   `form:"isOnline"`
	Seq    int32
}
