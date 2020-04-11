package http

import (
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/gin-gonic/gin"
)

func (s *Server) deviceRegister(c *gin.Context) {
	var device model.Device
	if err := c.BindJSON(&device); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err := s.logic.DeviceRegister(c, &device); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	var res = make(map[string]interface{}, 2)
	res["deviceId"] = device.ID
	res["goimKey"] = device.Key
	result(c, res, OK)
}

func (s *Server) deviceStatus(c *gin.Context) {
	var arg struct {
		Sn string `form:"sn"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	var (
		device *model.Device
		err    error
	)

	if device, err = s.logic.DeviceStatus(c, arg.Sn); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, device, OK)
}
