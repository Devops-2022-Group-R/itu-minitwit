package database

import (
	"errors"

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
	var latest LatestDTO
	err := rep.db.Last(&latest).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			latest = LatestDTO{Value: newLatest}
			return rep.db.Create(&latest).Error
		}

		return err
	}

	latest.Value = newLatest
	return rep.db.Save(&latest).Error
}

func (rep *GormLatestRepository) GetCurrent() (int, error) {
	var dto LatestDTO
	err := rep.db.Take(&dto).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, nil
	}

	return dto.Value, err
}
