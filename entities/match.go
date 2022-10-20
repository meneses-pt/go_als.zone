package entities

import "time"

type Match struct {
	ID       uint       `json:"id"`
	Slug     *string    `json:"slug"`
	Datetime *time.Time `json:"datetime"`
	Score    *string    `json:"score"`
}
