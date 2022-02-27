package database

type LatestDTO struct {
	Id    int `gorm:"primaryKey"`
	Value int `gorm:"not null"`
}

func (LatestDTO) TableName() string {
	return "latest"
}

type ILatestRepository interface {
	Migrate() error
	Set(newValue int) error
	GetCurrent() (int, error)
}
