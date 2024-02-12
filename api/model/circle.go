package model

import (
	"gorm.io/gorm"
	"time"
)

type Circle struct {
	UpdatedAt   time.Time          `json:"updatedAt" gorm:"autoUpdateTime;"`
	CreatedAt   time.Time          `json:"createdAt" gorm:"autoCreateTime;"`
	ValidUntil  *time.Time         `json:"validUntil"`
	CreatedFrom string             `json:"createdFrom" gorm:"type:varchar(50);not null"`
	ImageSrc    string             `json:"imageSrc" gorm:"type:text;not null;"`
	Description string             `json:"description" gorm:"type:varchar(1200);not null;"`
	Name        string             `json:"name" gorm:"type:varchar(40);not null;"`
	Votes       []*Vote            `json:"votes" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	Voters      []*CircleVoter     `json:"voters" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	Candidates  []*CircleCandidate `json:"candidate" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	ID          int64              `json:"id" gorm:"primary_key;index;"`
	Private     bool               `json:"private" gorm:"not null;default:false;"`
	Active      bool               `json:"active" gorm:"not null;default:true;"`
}

type CircleUriRequest struct {
	CircleID int64 `uri:"circleId"`
}

type CircleByUriRequest struct {
	Name string `uri:"name" validate:"lte=40"`
}

type CircleResponse struct {
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ValidUntil  *time.Time `json:"validUntil"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageSrc    string     `json:"imageSrc"`
	CreatedFrom string     `json:"createdFrom"`
	ID          int64      `json:"id"`
	Private     bool       `json:"private"`
	Active      bool       `json:"active"`
}

type CircleUpdateRequest struct {
	Name        *string                   `json:"name,omitempty" validate:"omitempty,gt=0,lte=40"`
	Description *string                   `json:"description,omitempty" validate:"omitempty,lte=1200"`
	ImageSrc    *string                   `json:"imageSrc,omitempty" validate:"omitempty,url"`
	Delete      *bool                     `json:"delete,omitempty" validate:"omitempty"`
	ValidUntil  *time.Time                `json:"validUntil,omitempty" validate:"omitempty"`
	Voters      []*CircleVoterRequest     `json:"voters,omitempty"`
	Candidates  []*CircleCandidateRequest `json:"candidates,omitempty"`
	ID          int64                     `json:"id" validate:"gt=0"`
}

type CircleCreateRequest struct {
	Description *string                   `json:"description,omitempty" validate:"omitempty,gt=0,lte=1200"`
	ImageSrc    *string                   `json:"imageSrc,omitempty" validate:"omitempty,url"`
	Private     *bool                     `json:"private,omitempty" validate:"omitempty"`
	ValidUntil  *time.Time                `json:"validUntil,omitempty" validate:"omitempty"`
	Name        string                    `json:"name" validate:"gt=0,lte=40"`
	Voters      []*CircleVoterRequest     `json:"voters,omitempty"`
	Candidates  []*CircleCandidateRequest `json:"candidates,omitempty"`
}

type CirclePaginated struct {
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ImageSrc        string    `json:"imageSrc"`
	ID              int64     `json:"id"`
	VotersCount     int64     `json:"votersCount"`
	CandidatesCount int64     `json:"candidatesCount"`
	Active          bool      `json:"active"`
}

type CirclePaginatedResponse struct {
	VotersCount     *int64 `json:"votersCount"`
	CandidatesCount *int64 `json:"candidatesCount"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageSrc        string `json:"imageSrc"`
	ID              int64  `json:"id"`
	Active          bool   `json:"active"`
}

// db hooks with checks

// AfterFind checks whether the circle validation is expired and set
// the circle inactive if so. If the update of the column failed,
// the query will fail.
func (circle *Circle) AfterFind(tx *gorm.DB) (err error) {
	// check if any validation time is set
	if circle.ValidUntil == nil {
		return
	}

	if isValidationTimeExpired(circle) {
		circle.Active = false
		err := tx.Model(circle).Update("active", false).Error

		if err != nil {
			return err
		}
	}

	return
}

func isValidationTimeExpired(
	circle *Circle,
) bool {
	if time.Now().After(*circle.ValidUntil) {
		return true
	}
	return false
}
