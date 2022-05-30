package model

import "time"

type Circle struct {
	ID int64 `json:"id" gorm:"primary_key;index;"`

	// Name of the circle
	Name string `json:"name" gorm:"type:varchar(40);not null;"`

	// Votes that the circle contains, 0 or more.
	Votes []*Vote `json:"votes" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`

	// Voters that the circle contains, 0 or more.
	Voters []*CircleVoter `json:"voters" gorm:"foreignKey:CircleRefer;constraint:OnDelete:CASCADE;"`

	// Private indicates if this Circle should be private and thus visible only to the users
	// that are eligible.
	Private bool `json:"private" gorm:"not null;default:false;"`

	// Active indicates if this circle is active, and it is possible to vote
	Active bool `json:"active" gorm:"not null;default:true;"`

	// CreatedFrom identity id that has created this circle
	CreatedFrom string `json:"createdFrom" gorm:"type:varchar(50);not null"`

	// ValidUntil if given, sets the time until this circle is valid and active
	ValidUntil *time.Time `json:"validUntil" gorm:""`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CircleUpdateInput struct {
	Name       *string             `json:"name"`
	Voters     []*CircleVoterInput `json:"voters"`
	Private    *bool               `json:"private"`
	Delete     *bool               `json:"delete"`
	ValidUntil *time.Time          `json:"validUntil"`
}
