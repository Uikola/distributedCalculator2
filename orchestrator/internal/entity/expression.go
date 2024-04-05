package entity

import "time"

type Expression struct {
	ID           string    `json:"id"`
	Expression   string    `json:"expression"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	CalculatedAt time.Time `json:"calculated_at"`
	CalculatedBy time.Time `json:"calculated_by"`
}
