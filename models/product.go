package models

import "time"

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Desc         string    `json:"desc"`
	Price        int       `json:"price"`
	Stock        int       `json:"stock"`
	CategoryID   *int      `json:"category_id"`
	CategoryName *string   `json:"category_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
