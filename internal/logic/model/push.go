package model

type PushKeyMessage struct {
	Op   int32    `form:"operation"`
	Keys []string `form:"keys"`
	Seq  int32    `form:"seq" binding:"required"`
}

type PushMidsMessage struct {
	Op   int32   `form:"operation"`
	Mids []int64 `form:"mids"`
	Seq  int32   `form:"seq" binding:"required"`
}

type PushRoomMessage struct {
	Op   int32  `form:"operation" binding:"required"`
	Seq  int32  `form:"seq" binding:"required"`
	Type string `form:"type" binding:"required"`
	Room string `form:"room" binding:"required"`
}

type PushAllMessage struct {
	Op    int32 `form:"operation" binding:"required"`
	Seq   int32 `form:"seq" binding:"required"`
	Speed int32 `form:"speed"`
}
