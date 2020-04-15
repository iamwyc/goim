package logic

import (
	"context"
	"sort"

	"github.com/Terry-Mao/goim/internal/logic/model"
)

var (
	_emptyTops = make([]*model.Top, 0)
)

// OnlineTop get the top online.
func (l *Logic) OnlineTop(c context.Context, arg *model.TopIn, n int) (tops []*model.Top, err error) {
	room := model.DecodePlatformAndSeriasRoomKey(arg.Platform, arg.Serias)
	for key, cnt := range l.roomCount {
		if room != "" && room == key {
			top := &model.Top{
				RoomID: key,
				Count:  cnt,
			}
			tops = append(tops, top)
			break
		} else if room == "" {
			top := &model.Top{
				RoomID: key,
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

	room := model.DecodePlatformAndSeriasRoomKey(arg.Platform, arg.Serias)
	res.Platform = arg.Platform
	res.Serias = arg.Serias
	if room != "" {
		res.Count = l.roomCount[room]
	}
	return
}

// OnlineTotal get all online.
func (l *Logic) OnlineTotal(c context.Context) (int64, int64, int) {
	registerCount, _ := l.DeviceCount(c)
	return l.totalIPs, l.totalConns, registerCount
}
