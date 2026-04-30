package models

import "time"

type Transaction struct {
	Title     string    `json:"title" bson:"title"`
	Amount    float64   `json:"amount" bson:"amount"`
	Type      string    `json:"type" bson:"type"`
	Source    string    `json:"source" bson:"source"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Email string `json:"email" bson:"email"`
}