package model

import "time"

type Voucher struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Code       string    `gorm:"unique;not null" json:"code"`
	Discount   float64   `json:"discount"`
	ExpiryDate time.Time `json:"expiry_date"`
	Status     string    `json:"status"` // e.g., "active", "inactive"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
