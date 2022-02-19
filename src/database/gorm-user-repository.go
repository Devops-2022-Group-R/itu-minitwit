package database

import (
	"errors"

	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db}
}

func (rep *GormUserRepository) Migrate() error {
	err := rep.db.AutoMigrate(&FollowDTO{})
	if err != nil {
		return err
	}

	return rep.db.AutoMigrate(&UserDTO{})
}

func (rep *GormUserRepository) Create(user models.User) error {
	dto := userDomaintoDto(user)

	return rep.db.Create(&dto).Error
}

func (rep *GormUserRepository) GetByID(id int64) (*models.User, error) {
	var dto UserDTO
	err := rep.db.First(&dto, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userDtoToDomain(dto), err
}

func (rep *GormUserRepository) GetByUsername(username string) (*models.User, error) {
	var dto UserDTO
	err := rep.db.Where("username = ?", username).First(&dto).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userDtoToDomain(dto), err
}

func (rep *GormUserRepository) Follow(whoId int64, whomId int64) error {
	dto := FollowDTO{
		WhoId:  whoId,
		WhomId: whomId,
	}

	return rep.db.Create(&dto).Error
}

func (rep *GormUserRepository) Unfollow(whoId int64, whomId int64) error {
	dto := FollowDTO{
		WhoId:  whoId,
		WhomId: whomId,
	}

	return rep.db.Where("who_id = ? AND whom_id = ?", whoId, whomId).Delete(&dto).Error
}

func (rep *GormUserRepository) IsFollowing(whoId int64, whomId int64) (bool, error) {
	var dto FollowDTO
	err := rep.db.Where("who_id = ? AND whom_id = ?", whoId, whomId).First(&dto).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, err
}

func userDtoToDomain(dto UserDTO) *models.User {
	return &models.User{
		UserId:       int64(dto.ID),
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: dto.PasswordHash,
	}
}

func userDomaintoDto(dto models.User) UserDTO {
	return UserDTO{
		Model: gorm.Model{
			ID: uint(dto.UserId),
		},
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: dto.PasswordHash,
	}
}