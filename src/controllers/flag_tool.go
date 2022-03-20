package controllers

import (
	"net/http"
	"strconv"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/database"
	"github.com/Devops-2022-Group-R/itu-minitwit/src/internal"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidMessageId = internal.NewBadRequestError("invalid message id")
)

func FlagMessageById(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)
	msgId, err := strconv.Atoi(c.Param("msgid"))
	if err != nil {
		internal.AbortWithError(c, ErrInvalidMessageId)
		return
	}
	message, err := messageRepository.FlagByMsgId(msgId)
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, message)
}

func GetAllMessages(c *gin.Context) {
	messageRepository := c.MustGet(MessageRepositoryKey).(database.IMessageRepository)
	messages, err := messageRepository.GetWithLimit(-1)
	if err != nil {
		internal.AbortWithError(c, internal.NewInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, messages)
}
