package model

// Online ip and room online.
type Online struct {
	Server    string           `json:"server"`
	RoomCount map[string]int32 `json:"room_count"`
	Updated   int64            `json:"updated"`
}

// Top top sorted.
type Top struct {
	RoomID string `json:"room_id"`
	Count  int32  `json:"count"`
}

// TopIn top sorted.
type TopIn struct {
	Platform int32 `form:"platform"`
	Serias   int32 `form:"serias"`
	Limit    int   `form:"limit" binding:"required"`
}

// OnlineRoom top sorted.
type OnlineRoom struct {
	Platform int32 `form:"platform" binding:"required"`
	Serias   int32 `form:"serias" binding:"required"`
}

// OnlineRoomOutVO top sorted.
type OnlineRoomOutVO struct {
	Platform int32
	Serias   int32
	Count    int32
}
