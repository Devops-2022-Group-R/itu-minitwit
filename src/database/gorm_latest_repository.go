package database

import (
	"gorm.io/gorm"
)

type GormLatestRepository struct {
	db *gorm.DB
}

func NewGormLatestRepository(db *gorm.DB) *GormLatestRepository {
	return &GormLatestRepository{db}
}

func (rep *GormLatestRepository) Migrate() error {
	return rep.db.AutoMigrate(&LatestDTO{})
}

func (rep *GormLatestRepository) Set(newLatest int) error {
	var latestCount int64
	err := rep.db.Model(&LatestDTO{}).Count(&latestCount).Error
	if err != nil {
		return err
	}

	dto := LatestDTO{Id: 1, Value: newLatest}
	if latestCount == 0 {
		err = rep.db.Create(&dto).Error
	} else {
		err = rep.db.Select("*").Updates(dto).Error
	}

	return err
}

func (rep *GormLatestRepository) GetCurrent() (int, error) {
	var dto LatestDTO
	err := rep.db.Take(&dto).Error
	return dto.Value, err
}
