package user

import (
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAllUsers() ([]User, error) {
	var users []User
	err := s.db.Find(&users).Error
	return users, err
}
