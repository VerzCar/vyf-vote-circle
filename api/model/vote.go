package model

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Vote struct {
	ID           int64       `json:"id" gorm:"primary_key;index;"`
	VoterRefer   int64       `json:"voterRefer"`
	Voter        CircleVoter `json:"voter" gorm:"foreignKey:VoterRefer;constraint:OnDelete:RESTRICT;"`
	ElectedRefer int64       `json:"electedRefer"`
	Elected      CircleVoter `json:"elected" gorm:"foreignKey:ElectedRefer;constraint:OnDelete:RESTRICT;"`
	CircleID     int64       `json:"circleId" gorm:"not null;"`
	Circle       *Circle     `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	CircleRefer  *int64      `json:"circleRefer"`
	CreatedAt    time.Time   `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt    time.Time   `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type VoteCreateRequest struct {
	CircleID int64  `json:"circleId" validate:"gt=0"`
	Elected  string `json:"elected" validate:"gt=0,lte=50"`
}
