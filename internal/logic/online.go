package logic

import (
	"context"
	"sort"
	"strings"

	"github.com/Terry-Mao/goim/internal/logic/model"
)

var (
	_emptyTops = make([]*model.Top, 0)
)

// OnlineTop get the top online.
func (l *Logic) OnlineTop(c context.Context, arg *model.TopIn, n int) (tops []*model.Top, err error) {
	rooms := model.DecodePlatformAndSeriasRoomKey(arg.Platform, arg.Serias)
	size := len(rooms)
	if size == 0 {
		return
	}
	for key, cnt := range l.roomCount {
		var roomKey string
		if arg.Platform > 0 && strings.HasPrefix(key, "p") {
			roomKey = model.EncodePlatformRoomKey(arg.Platform)
		} else if arg.Serias > 0 && strings.HasPrefix(key, "s") {
			roomKey = model.EncodeSeriasRoomKey(arg.Serias)
		}
		if roomKey != "" {
			top := &model.Top{
				RoomID: roomKey,
				Count:  cnt,
			}
			tops = append(tops, top)
		}
	}
	sort.Slice(tops, func(i, j int) bool {
		return tops[i].Count > tops[j].Count
	})
	if len(tops) > n {
		tops = tops[:n]
	}
	if len(tops) == 0 {
		tops = _emptyTops
	}
	return
}

// OnlineRoom get rooms online.
func (l *Logic) OnlineRoom(c context.Context, arg *model.OnlineRoom) (res model.OnlineRoomOutVO, err error) {
	if arg.Platform > 0 {
		res.Platform = l.roomCount[model.EncodePlatformRoomKey(arg.Platform)]
	}
	if arg.Serias > 0 {
		res.Serias = l.roomCount[model.EncodePlatformRoomKey(arg.Serias)]
	}
	return
}

// OnlineTotal get all online.
func (l *Logic) OnlineTotal(c context.Context) (int64, int64, int) {
	registerCount, _ := l.DeviceCount(c)
	return l.totalIPs, l.totalConns, registerCount
}
