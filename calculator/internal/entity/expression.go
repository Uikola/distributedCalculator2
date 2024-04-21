package entity

import "time"

type Expression struct {
	ID           uint             `json:"id"`
	Expression   string           `json:"expression"`
	Result       string           `json:"result"`
	Status       ExpressionStatus `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	CalculatedAt time.Time        `json:"calculated_at"`
	CalculatedBy uint             `json:"calculated_by"`
	OwnerID      uint             `json:"owner_id"`
}

type ExpressionStatus string

const (
	Error      ExpressionStatus = "error"
	InProgress ExpressionStatus = "in_progress"
	OK         ExpressionStatus = "ok"
)
