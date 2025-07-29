package model

type Game struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	Products    []Product `json:"products" gorm:"foreignKey:GameID"`
}
