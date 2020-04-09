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
