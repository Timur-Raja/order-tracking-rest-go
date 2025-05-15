package services

import (
	"errors"

	"github.com/timur-raja/order-tracking-rest-go/internal/models"
)

// CreateOrder creates a new order and returns the created order.
func CreateOrder(order models.Order) (models.Order, error) {
	// Business logic to create an order
	if order.ID == 0 {
		return models.Order{}, errors.New("order ID cannot be empty")
	}
	// Assume order is created successfully
	return order, nil
}

// GetUser retrieves a user by ID.
func GetUser(userID int) (models.User, error) {
	// Business logic to retrieve a user
	if userID == 0 {
		return models.User{}, errors.New("user ID cannot be empty")
	}
	// Assume we found the user
	return models.User{ID: userID}, nil
}
