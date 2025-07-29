package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GameID      uint      `json:"game_id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
