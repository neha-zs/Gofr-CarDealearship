package models

import "github.com/google/uuid"

type Engine struct {
	EngineID     uuid.UUID `json:"id,omitempty"`
	Displacement int       `json:"displacement,omitempty"`
	Cylinders    int       `json:"cylinders,omitempty"`
	Range        int       `json:"range,omitempty"`
}
