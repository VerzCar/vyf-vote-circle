package factory

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"time"
)

var LocaleGermany = &model.Locale{
	ID:           83,
	Locale:       "Germany",
	LanguageCode: "de",
	LcidString:   "de-de",
	CreatedAt:    time.Time{},
	UpdatedAt:    time.Time{},
}
