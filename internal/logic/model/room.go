package model

import (
	"fmt"
)

const (
	platformRoomFmt       = "p://%d"
	seriasRoomFmt         = "s://%d"
	platfromSeriasRoomFmt = "ps://p%ds%d"
)

// DecodePlatformAndSeriasRoomKey encode a room key.
func DecodePlatformAndSeriasRoomKey(platform int32, serias int32) (room string) {

	if platform > 0 && serias > 0 {
		room = platformAndSeriasRoom(platform, serias)
	} else if platform > 0 {
		room = platformRoom(platform)
	} else if serias > 0 {
		room = seriasRoom(serias)
	}

	return room
}

// EncodePlatformAndSeriasRoomKey encode a room key.
func EncodePlatformAndSeriasRoomKey(platform int32, serias int32) string {
	room := ""
	if platform > 0 {
		room = room + platformRoom(platform) + "@"
	}
	if serias > 0 {
		room = room + seriasRoom(serias) + "@"
	}
	if platform > 0 && serias > 0 {
		room = room + platformAndSeriasRoom(platform, serias) + "@"
	}
	if room != "" {
		room = room[0 : len(room)-1]
	}
	return room
}

func platformRoom(platform int32) string {
	if platform > 0 {
		return fmt.Sprintf(platformRoomFmt, platform)
	}
	return ""
}

func seriasRoom(serias int32) string {
	if serias > 0 {
		return fmt.Sprintf(seriasRoomFmt, serias)
	}
	return ""
}

func platformAndSeriasRoom(platform int32, serias int32) string {
	if platform > 0 && serias > 0 {
		return fmt.Sprintf(platfromSeriasRoomFmt, platform, serias)
	}
	return ""
}
