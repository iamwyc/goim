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
