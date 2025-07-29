package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	Email         string    `gorm:"unique" json:"email"`
	Password      string    `json:"-"`
	PhoneNumber   string    `json:"phone_number"`
	Role          string    `json:"role"`
	RememberToken *string   `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
