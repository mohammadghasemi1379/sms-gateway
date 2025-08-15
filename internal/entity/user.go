package entity

import "time"

type User struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	Credit    int64     `json:"credit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
