package logic

import (
	"context"

	"github.com/Terry-Mao/goim/internal/logic/model"
)

// DeviceRegister device register
func (l *Logic) DeviceRegister(c context.Context, arg *model.Device) error {
	return l.dao.DeviceRegister(arg)
}

// DeviceCount device register
func (l *Logic) DeviceCount(c context.Context) (int, error) {
	return l.dao.DeviceCount()
}

// DeviceStatus device status
func (l *Logic) DeviceStatus(c context.Context, sn string) (*model.Device, error) {
	d, e := l.dao.GetDeviceBySn(sn)
	if e != nil {
		return d, e
	}
	res, e := l.dao.ServersByKeys(c, []string{sn})
	if e == nil {
		d.Online = len(res) > 0
	}
	return d, e
}
