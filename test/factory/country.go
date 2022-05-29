package factory

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"time"
)

var CountryGermany = &model.Country{
	ID:            84,
	Name:          "Germany",
	Alpha2:        "DE",
	Alpha3:        "DEU",
	ContinentCode: "EU",
	Number:        "276",
	FullName:      "Federal Republic of Germany",
	CreatedAt:     time.Time{},
	UpdatedAt:     time.Time{},
}
