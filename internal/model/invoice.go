package model

import "time"

type Invoice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `json:"order_id"`
	Order       Order     `json:"order" gorm:"foreignKey:OrderID"`
	InvoiceCode string    `json:"invoice_code"` // INV-xxxxxx
	GeneratedAt time.Time `json:"generated_at"`
}
