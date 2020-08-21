package http

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/golang/glog"
	"github.com/google/uuid"

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
	var msgID int64
	if msgID, err = s.logic.PushSnList(context.TODO(), &arg, msg); err != nil {
		result(c, nil, RequestErr)
		return
	}
	result(c, msgID, OK)
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
	var msgID int64
	if msgID, err = s.logic.PushMidList(context.TODO(), &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, msgID, OK)
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
	var id int64
	if id, err = s.logic.PushRoom(c, &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, id, OK)
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
	var msgID int64
	if msgID, err = s.logic.PushAll(c, &arg, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, msgID, OK)
}

func (s *Server) pushBySnFile(c *gin.Context) {
	var arg model.PushMessageIDParam
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	glog.Infof("messageId:%d", arg)
	form, _ := c.MultipartForm()
	files := form.File["files"]
	var fileList []string
	for _, file := range files {
		filePath := s.logic.GetFileDir() + strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")) + ".txt"
		err := c.SaveUploadedFile(file, filePath)
		if err != nil {
			glog.Error("save error:", err)
		} else {
			fileList = append(fileList, filePath)
		}
	}
	go s.logic.PushSnFils(arg.MessageID, fileList)
	result(c, nil, OK)
}

func (s *Server) pushStatus(c *gin.Context) {
	var arg model.PushMessageIDParam
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	glog.Infof("messageId:%d", arg)
	count, err := s.logic.MessageRedisStats(c, arg.MessageID)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	result(c, count, OK)
}
