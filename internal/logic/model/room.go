package model

import (
	"fmt"
	"net/url"
)
const(
	platfromRoom ="p://%d"
	seriasRoom ="s://%d"
)
// EncodeRoomKey encode a room key.
func EncodeRoomKey(typ string, room string) string {
	return fmt.Sprintf("%s://%s", typ, room)
}

// EncodePlatformAndSeriasRoomKey encode a room key.
func EncodePlatformAndSeriasRoomKey(platform int32, serias int32) string {
	room := ""
	if platform > 0 {
		room = room + EncodePlatformRoomKey(platform) + "@"
	}
	if serias > 0 {
		room = room + EncodeSeriasRoomKey(serias) + "@"
	}
	if room != "" {
		room = room[0 : len(room)-1]
	}
	return room
}

// EncodePlatformRoomKey encode a room key.
func EncodePlatformRoomKey(platform int32) string {
	return fmt.Sprintf(platfromRoom, platform)
}

// EncodeSeriasRoomKey encode a room key.
func EncodeSeriasRoomKey(serias int32) string {
	return fmt.Sprintf(seriasRoom, serias)
}

// DecodeRoomKey decode room key.
func DecodeRoomKey(key string) (string, string, error) {
	u, err := url.Parse(key)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}
