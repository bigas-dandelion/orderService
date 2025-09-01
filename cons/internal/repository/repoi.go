package repository

import "l0/cons/internal/models"

type IRepository interface {
	StoreOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
}
