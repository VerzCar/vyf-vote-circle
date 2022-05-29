package model

import "time"

type CircleVoter struct {
	ID int64 `json:"id" gorm:"primary_key;index;"`

	// Voter that is eligible to vote.
	// This must be a user identity id that should vote for this circle.
	Voter string `json:"voter" gorm:"type:varchar(50);not null"`

	// Committed has committed to vote.
	Committed bool `json:"committed" gorm:"default:false;not null"`

	// Rejected has rejected to vote.
	Rejected bool `json:"rejected" gorm:"default:false;not null"`

	CircleID    int64     `json:"circleId" gorm:"not null;"`
	Circle      *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	CircleRefer *int64    `json:"circleRefer"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}
