package routes

import (
	"pasarmalam/models"
)

// Aliases to keep routes.go free of model import list
type (
	databaseUser     = models.User
	databaseProduct  = models.Product
	databaseCart     = models.Cart
	databaseCartItem = models.CartItem
	databaseOrder    = models.Order
	databaseOrderItem = models.OrderItem
)
