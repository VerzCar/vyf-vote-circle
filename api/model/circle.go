package model

import (
	"database/sql/driver"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"gorm.io/gorm"
	"time"
)

type Circle struct {
	UpdatedAt   time.Time          `json:"updatedAt" gorm:"autoUpdateTime;"`
	CreatedAt   time.Time          `json:"createdAt" gorm:"autoCreateTime;"`
	ValidFrom   time.Time          `json:"validFrom"`
	ValidUntil  *time.Time         `json:"validUntil"`
	CreatedFrom string             `json:"createdFrom" gorm:"type:varchar(50);not null"`
	ImageSrc    string             `json:"imageSrc" gorm:"type:text;not null;"`
	Description string             `json:"description" gorm:"type:varchar(1200);not null;"`
	Name        string             `json:"name" gorm:"type:varchar(40);not null;"`
	Votes       []*Vote            `json:"votes" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	Voters      []*CircleVoter     `json:"voters" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	Candidates  []*CircleCandidate `json:"candidate" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	Stage       CircleStage        `json:"stage" gorm:"type:circleStage;not null;default:COLD"`
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
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	ValidFrom   time.Time   `json:"validFrom"`
	ValidUntil  *time.Time  `json:"validUntil"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ImageSrc    string      `json:"imageSrc"`
	CreatedFrom string      `json:"createdFrom"`
	Stage       CircleStage `json:"stage"`
	ID          int64       `json:"id"`
	Private     bool        `json:"private"`
	Active      bool        `json:"active"`
}

type CircleUpdateRequest struct {
	Name        *string                   `json:"name,omitempty" validate:"omitempty,gt=0,lte=40"`
	Description *string                   `json:"description,omitempty" validate:"omitempty,lte=1200"`
	ImageSrc    *string                   `json:"imageSrc,omitempty" validate:"omitempty,url"`
	Delete      *bool                     `json:"delete,omitempty" validate:"omitempty"`
	ValidUntil  *time.Time                `json:"validUntil,omitempty" validate:"omitempty"`
	ValidFrom   *time.Time                `json:"ValidFrom,omitempty" validate:"omitempty"`
	Voters      []*CircleVoterRequest     `json:"voters,omitempty"`
	Candidates  []*CircleCandidateRequest `json:"candidates,omitempty"`
	ID          int64                     `json:"id" validate:"gt=0"`
}

type CircleCreateRequest struct {
	Description *string                   `json:"description,omitempty" validate:"omitempty,gt=0,lte=1200"`
	ImageSrc    *string                   `json:"imageSrc,omitempty" validate:"omitempty,url"`
	Private     *bool                     `json:"private,omitempty" validate:"omitempty"`
	ValidUntil  *time.Time                `json:"validUntil,omitempty" validate:"omitempty"`
	ValidFrom   *time.Time                `json:"ValidFrom,omitempty" validate:"omitempty"`
	Name        string                    `json:"name" validate:"gt=0,lte=40"`
	Voters      []*CircleVoterRequest     `json:"voters,omitempty"`
	Candidates  []*CircleCandidateRequest `json:"candidates,omitempty"`
}

type CirclePaginated struct {
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	ImageSrc        string      `json:"imageSrc"`
	Stage           CircleStage `json:"stage"`
	ID              int64       `json:"id"`
	VotersCount     int64       `json:"votersCount"`
	CandidatesCount int64       `json:"candidatesCount"`
	Active          bool        `json:"active"`
}

type CirclePaginatedResponse struct {
	VotersCount     *int64      `json:"votersCount"`
	CandidatesCount *int64      `json:"candidatesCount"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	ImageSrc        string      `json:"imageSrc"`
	Stage           CircleStage `json:"stage"`
	ID              int64       `json:"id"`
	Active          bool        `json:"active"`
}

type CircleStage string

const (
	CircleStageCold   CircleStage = "COLD"
	CircleStageHot    CircleStage = "HOT"
	CircleStageClosed CircleStage = "CLOSED"
)

func (e *CircleStage) Scan(value interface{}) error {
	*e = CircleStage(value.(string))
	return nil
}

func (e CircleStage) Value() (driver.Value, error) {
	return string(e), nil
}

func (e CircleStage) IsValid() bool {
	switch e {
	case CircleStageCold, CircleStageHot, CircleStageClosed:
		return true
	}
	return false
}

func (e CircleStage) String() string {
	return string(e)
}

// db hooks with checks ++++++++++++++++++++++++++++

func (circle *Circle) AfterFind(tx *gorm.DB) error {
	if !circle.Active || circle.Stage == CircleStageClosed {
		return nil
	}

	return updateCircleStage(tx, circle)
}

func (circle *Circle) BeforeCreate(tx *gorm.DB) error {
	return updateCircleStage(tx, circle)
}

func (circle *Circle) BeforeUpdate(tx *gorm.DB) error {
	return updateCircleStage(tx, circle)
}

func updateCircleStage(
	tx *gorm.DB,
	circle *Circle,
) error {
	currentTime := time.Now().Truncate(24 * time.Hour)
	validFromTruncatedTime := circle.ValidFrom.Truncate(24 * time.Hour)
	validUntilTime := *circle.ValidUntil
	validUntilTruncatedTime := validUntilTime.Truncate(24 * time.Hour)

	// check if current time is between range of circle
	// if so, set it to hot stage
	if circle.ValidUntil != nil {
		if utils.IsTimeBetween(currentTime, validFromTruncatedTime, validUntilTruncatedTime) {
			err := tx.Model(circle).Update("stage", CircleStageHot).Error

			if err != nil {
				return err
			}
			return nil
		}

		// check if current time is after valid until of circle
		// if so, set it to closed stage
		if currentTime.After(validUntilTruncatedTime) {
			err := tx.Model(circle).Update("stage", CircleStageClosed).Error

			if err != nil {
				return err
			}
			return nil
		}

		if currentTime.Before(validFromTruncatedTime) {
			err := tx.Model(circle).Update("stage", CircleStageCold).Error

			if err != nil {
				return err
			}
			return nil
		}
	}

	// check if current time is after valid from of circle
	// if so, set it to hot stage
	if currentTime.After(validFromTruncatedTime) {
		err := tx.Model(circle).Update("stage", CircleStageHot).Error

		if err != nil {
			return err
		}
		return nil
	}

	if currentTime.Before(validFromTruncatedTime) {
		err := tx.Model(circle).Update("stage", CircleStageCold).Error

		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
