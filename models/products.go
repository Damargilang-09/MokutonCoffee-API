package models

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`
	Name        string         `json:"name"`
	Image       string         `json:"image"`
	ImageId     string         `json:"imageid"`
	Description string         `json:"description"`
	Stock       int            `json:"stock"`
	Price       int64          `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
