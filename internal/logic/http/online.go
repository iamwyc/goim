package http

import (
	"context"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/gin-gonic/gin"
)

func (s *Server) onlineTop(c *gin.Context) {
	var arg model.TopIn
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	res, err := s.logic.OnlineTop(c, &arg, arg.Limit)
	if err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, res, OK)
}

func (s *Server) onlineRoom(c *gin.Context) {
	var arg model.OnlineRoom
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	res, err := s.logic.OnlineRoom(c, &arg)
	if err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, res, OK)
}

func (s *Server) onlineTotal(c *gin.Context) {
	ipCount, connCount, rc := s.logic.OnlineTotal(context.TODO())
	res := map[string]interface{}{
		"ip_count":       ipCount,
		"conn_count":     connCount,
		"register_count": rc,
	}
	result(c, res, OK)
}
