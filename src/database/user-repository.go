package database

import (
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"gorm.io/gorm"
)

const userTable = "user"

type UserDTO struct {
	gorm.Model

	Username     string `gorm:"not null;index"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`

	Followed  []FollowDTO `gorm:"foreignkey:WhoId"`
	Followers []FollowDTO `gorm:"foreignkey:WhomId"`
}

func (UserDTO) TableName() string {
	return userTable
}

type FollowDTO struct {
	WhoId int64   `gorm:"not null;index;index:who_id_whom_id"`
	Who   UserDTO `gorm:"primaryKey;foreignkey:WhoId"`

	WhomId int64   `gorm:"not null;index:who_id_whom_id"`
	Whom   UserDTO `gorm:"primaryKey;foreignkey:WhomId"`
}

func (FollowDTO) TableName() string {
	return "follower"
}

type IUserRepository interface {
	Migrate() error
	Create(users models.User) error
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	NumUsers() (int64, error)

	Follow(whoId int64, whomId int64) error
	Unfollow(whoId int64, whomId int64) error
	IsFollowing(whoId int64, whomId int64) (bool, error)
	AllFollowed(whoId int64) ([]*models.User, error)
}
