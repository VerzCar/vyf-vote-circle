package model

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Vote struct {
	Voter          *CircleVoter     `json:"voter" gorm:"foreignKey:VoterRefer;constraint:OnDelete:RESTRICT;"`
	Candidate      *CircleCandidate `json:"candidate" gorm:"foreignKey:CandidateRefer;constraint:OnDelete:RESTRICT;"`
	CreatedAt      time.Time        `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt      time.Time        `json:"updatedAt" gorm:"autoUpdateTime;"`
	Circle         *Circle          `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	CircleRefer    *int64           `json:"circleRefer"`
	ID             int64            `json:"id" gorm:"primary_key;index;"`
	VoterRefer     int64            `json:"voterRefer"`
	CandidateRefer int64            `json:"candidateRefer"`
	CircleID       int64            `json:"circleId" gorm:"not null;"`
}

type VoteCreateRequest struct {
	CandidateID string `json:"candidateId" validate:"gt=0,lte=50"`
}
