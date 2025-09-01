package services

import (
	"l0/cons/internal/models"
	"l0/cons/internal/repository"
)

type Service struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) StoreOrder(order *models.Order) error {
	return s.repo.StoreOrder(order)
}

func (s *Service) GetOrder(orderUID string) (*models.Order, error) {
	res, err := s.repo.GetOrder(orderUID)
	return res, err
}