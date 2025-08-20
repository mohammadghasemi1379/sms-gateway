package entity

import "time"

type User struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	PhoneNumber string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	Credit      int64     `json:"credit"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
