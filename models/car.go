package models

import "github.com/google/uuid"

type Car struct {
	ID       uuid.UUID `json:"ID,omitempty"`
	Engine   Engine    `json:"Engine,omitempty"`
	Name     string    `json:"Name"`
	Year     int       `json:"Year"`
	Brand    string    `json:"Brand"`
	FuelType string    `json:"FuelType"`
}
