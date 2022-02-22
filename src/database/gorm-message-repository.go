package database

import (
	"github.com/Devops-2022-Group-R/itu-minitwit/src/models"
	"gorm.io/gorm"
)

type GormMessageRepository struct {
	db *gorm.DB
}

func NewGormMessageRepository(db *gorm.DB) *GormMessageRepository {
	return &GormMessageRepository{db}
}

func (rep *GormMessageRepository) Migrate() error {
	return rep.db.AutoMigrate(&MessageDTO{})
}

func (rep *GormMessageRepository) Create(message models.Message) error {
	dto := MessageDTO{
		Author:  userDomaintoDto(message.Author),
		Text:    message.Text,
		PubDate: message.PubDate,
		Flagged: message.Flagged,
	}

	return rep.db.Create(&dto).Error
}

func (rep *GormMessageRepository) GetWithLimit(limit int) ([]models.Message, error) {
	var dtos []MessageDTO
	err := rep.db.Order("pub_date desc").Preload("Author").Limit(limit).Find(&dtos).Error

	messages := make([]models.Message, len(dtos))
	for i, dto := range dtos {
		messages[i] = messageDtoToDomain(dto)
	}

	return messages, err
}

func (rep *GormMessageRepository) GetByUserId(userId int64, limit int) ([]models.Message, error) {
	var dtos []MessageDTO
	err := rep.db.Order("pub_date desc").Preload("Author").Where("author_id = ?", userId).Limit(limit).Find(&dtos).Error

	messages := make([]models.Message, len(dtos))
	for i, dto := range dtos {
		messages[i] = messageDtoToDomain(dto)
	}

	return messages, err
}

func (rep *GormMessageRepository) GetByUserAndItsFollowers(userId int64, limit int) ([]models.Message, error) {
	var dtos []MessageDTO

	subQuery := rep.db.Table("follower").Select("whom_id").Where("who_id = ?", userId)
	err := rep.db.Order("pub_date desc").Preload("Author").Where("flagged IS FALSE AND (author_id = ? OR author_id = (?))", userId, subQuery).Limit(limit).Find(&dtos).Error

	messages := make([]models.Message, len(dtos))
	for i, dto := range dtos {
		messages[i] = messageDtoToDomain(dto)
	}

	return messages, err
}

func messageDtoToDomain(dto MessageDTO) models.Message {
	return models.Message{
		Author: models.User{
			UserId:   int64(dto.Author.ID),
			Username: dto.Author.Username,
			Email:    dto.Author.Email,
		},
		Text:    dto.Text,
		PubDate: dto.PubDate,
		Flagged: dto.Flagged,
	}
}
