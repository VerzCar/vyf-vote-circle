package model

import "time"

type UserOption struct {
	UpdatedAt            time.Time           `json:"updatedAt" gorm:"autoUpdateTime;"`
	CreatedAt            time.Time           `json:"createdAt" gorm:"autoCreateTime;"`
	IdentityID           string              `json:"identityId" gorm:"type:varchar(50);not null"`
	Package              SubscriptionPackage `json:"package" gorm:"type:subscriptionPackage;not null;default:S"`
	ID                   int64               `json:"id" gorm:"primary_key;index;"`
	MaxCircles           int                 `json:"maxCircles" gorm:"not null"`
	MaxVoters            int                 `json:"maxVoters" gorm:"not null"`
	MaxCandidates        int                 `json:"maxCandidates" gorm:"not null"`
	MaxPrivateVoters     int                 `json:"maxPrivateVoters" gorm:"not null"`
	MaxPrivateCandidates int                 `json:"maxPrivateCandidates" gorm:"not null"`
}

type UserOptionResponse struct {
	MaxCircles    int                       `json:"maxCircles"`
	MaxVoters     int                       `json:"maxVoters"`
	MaxCandidates int                       `json:"maxCandidates"`
	PrivateOption UserPrivateOptionResponse `json:"privateOption"`
}

type UserPrivateOptionResponse struct {
	MaxVoters     int `json:"maxVoters"`
	MaxCandidates int `json:"maxCandidates"`
}
