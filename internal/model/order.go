package model

import "time"

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	ProductID     uint      `json:"product_id"`
	Product       Product   `json:"product" gorm:"foreignKey:ProductID"`
	GameAccountID string    `json:"game_account_id"`
	Quantity      int       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"` // ← TRX-xxxxxx (kita generate)
	MidtransID    string    `json:"midtrans_id"`    // ← dari Midtrans (optional)
	PaymentProof  *string   `json:"payment_proof"`
	VoucherID     *uint     `json:"voucher_id"`
	Voucher       *Voucher  `json:"voucher"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
