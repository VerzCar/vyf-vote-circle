package model

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Vote struct {
	ID          int64     `json:"id" gorm:"primary_key;index;"`
	Voter       string    `json:"voter" gorm:"type:varchar(50);not null"`
	Elected     string    `json:"elected" gorm:"type:varchar(50);not null"`
	CircleID    int64     `json:"circleId" gorm:"not null;"`
	Circle      *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	CircleRefer *int64    `json:"circleRefer"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}
