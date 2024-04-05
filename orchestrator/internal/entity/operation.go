package entity

import "time"

type Operation struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}
