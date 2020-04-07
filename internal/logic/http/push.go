package http

import (
	"context"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func (s *Server) pushSnList(c *gin.Context) {
	var arg model.PushKeyMessage
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	// read message
	arg.Op = model.DefaultOperation
	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.PushSnList(context.TODO(), &arg, msg); err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, nil, OK)
}

func (s *Server) pushMidList(c *gin.Context) {

	var arg model.PushMidsMessage
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	// read message
	arg.Op = model.DefaultOperation
	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.PushMidList(context.TODO(), &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}

func (s *Server) pushRoom(c *gin.Context) {
	var arg model.PushRoomMessage
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	// read message
	arg.Op = model.DefaultOperation
	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.PushRoom(c, &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}

func (s *Server) pushAll(c *gin.Context) {
	var arg model.PushAllMessage
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	arg.Op = model.DefaultOperation
	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if err = s.logic.PushAll(c, &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}
