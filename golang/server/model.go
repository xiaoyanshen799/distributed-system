package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetModel represents the pet entity in the database
type PetModel struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Gender  string
	Age     int32
	Breed   string
	Picture string
}

// BeforeCreate hook to set UUID before creating a new record
func (p *PetModel) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
