package database

import (
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"gorm.io/gorm"
)

type MessageDTO struct {
	gorm.Model

	AuthorId int     `gorm:"not null"`
	Author   UserDTO `gorm:"foreignkey:AuthorId"`
	Text     string  `gorm:"not null"`
	PubDate  int64   `gorm:"not null;index:,sort:desc"`
	Flagged  bool    `gorm:"not null"`
}

func (MessageDTO) TableName() string {
	return "message"
}

type IMessageRepository interface {
	Migrate() error
	Create(message models.Message) error
	GetWithLimit(limit int) ([]models.Message, error)
	GetByUserId(userId int64, limit int) ([]models.Message, error)
	GetByUserAndItsFollowers(userId int64, limit int) ([]models.Message, error)
	FlagByMsgId(msgId int) (models.Message, error)
}
