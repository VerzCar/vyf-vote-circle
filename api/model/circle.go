package model

import (
	"gorm.io/gorm"
	"time"
)

type Circle struct {
	ID int64 `json:"id" gorm:"primary_key;index;"`
	// Name of the circle
	Name string `json:"name" gorm:"type:varchar(40);not null;"`
	// Description of the circle
	Description string `json:"description" gorm:"type:varchar(1200);not null;"`
	// ImageSrc of the circle
	ImageSrc string `json:"imageSrc" gorm:"type:text;not null;"`
	// Votes that the circle contains, 0 or more.
	Votes []*Vote `json:"votes" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	// Voters that the circle contains, 0 or more.
	Voters []*CircleVoter `json:"voters" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`
	// Private indicates if this Circle should be private and thus visible only to the users
	// that are eligible.
	Private bool `json:"private" gorm:"not null;default:false;"`
	// Active indicates if this circle is active, and it is possible to vote.
	// Otherwise, it is closed
	Active bool `json:"active" gorm:"not null;default:true;"`
	// CreatedFrom identity id that has created this circle
	CreatedFrom string `json:"createdFrom" gorm:"type:varchar(50);not null"`
	// ValidUntil if given, sets the time until this circle is valid and active
	ValidUntil *time.Time `json:"validUntil"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CircleRequest struct {
	ID int64 `json:"id"`
}

type CircleResponse struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ImageSrc    string                 `json:"imageSrc"`
	Voters      []*CircleVoterResponse `json:"voters"`
	Private     bool                   `json:"private"`
	Active      bool                   `json:"active"`
	CreatedFrom string                 `json:"createdFrom"`
	ValidUntil  *time.Time             `json:"validUntil"`
}

type CircleUpdateInput struct {
	Name        *string             `json:"name" validate:"omitempty,gt=0,lte=40"`
	Description *string             `json:"description" validate:"omitempty,gt=0,lte=1200"`
	ImageSrc    *string             `json:"imageSrc" validate:"omitempty,url"`
	Voters      []*CircleVoterInput `json:"voters"`
	Private     *bool               `json:"private" validate:"omitempty"`
	Delete      *bool               `json:"delete" validate:"omitempty"`
	ValidUntil  *time.Time          `json:"validUntil" validate:"omitempty"`
}

type CircleCreateInput struct {
	Name        string              `json:"name" validate:"gt=0,lte=40"`
	Description *string             `json:"description" validate:"omitempty,gt=0,lte=1200"`
	ImageSrc    *string             `json:"imageSrc" validate:"omitempty,url"`
	Voters      []*CircleVoterInput `json:"voters"`
	Private     *bool               `json:"private" validate:"omitempty"`
	ValidUntil  *time.Time          `json:"validUntil" validate:"omitempty"`
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
