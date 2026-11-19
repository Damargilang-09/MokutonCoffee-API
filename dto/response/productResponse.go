package response

import "github.com/google/uuid"

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	Price       int64     `json:"price"`
}
