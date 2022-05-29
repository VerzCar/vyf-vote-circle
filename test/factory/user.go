package factory

import (
	"github.com/golang-jwt/jwt/v4"
	"gitlab.vecomentman.com/libs/sso"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

type User struct {
	Martin           model.User
	Albert           model.User
	Marco            model.User
	MarcosPizzaUsers []model.User
	Sarah            model.User
}

type InputUser struct {
	Elon model.UserCreateInput
}

type SsoUser struct {
	Elon   sso.SsoClaims
	Martin sso.SsoClaims
}

type CreationUser struct {
	Elon model.User
}

type UpdateInputUser struct {
	Martin model.UserUpdateInput
}

var martinEmail = "dev.martin@vecomentman.de"
var martin = model.User{
	ID:         1,
	IdentityID: "ee7292f9-6e2c-4825-9d1e-ef6776e954b4",
	FirstName:  "Martin",
	LastName:   "Hammer",
	Gender:     model.GenderMan,
	Locale:     LocaleGermany,
	Address: &model.Address{
		Address:    "Badestr. 4",
		City:       "Laufen",
		PostalCode: "83410",
		Country:    CountryGermany,
	},
}

var genderUpdate = model.GenderWomen
var martinUpdateInput = model.UserUpdateInput{
	FirstName: stringP("Morten"),
	LastName:  stringP("LastMorten"),
	Gender:    &genderUpdate,
	Locale:    &LocaleGermany.LcidString,
	AvatarURL: stringP("https://media.vecomentman.de/e4893f00o-4384834-2332/avatar"),
	Address:   &martinUpdateAddressInput,
	Contact:   &martinUpdateContactInput,
}

var martinUpdateAddressInput = model.AddressInput{
	Address:          "Badestr. 23",
	City:             "Genfo",
	PostalCode:       "83410",
	CountryAlphaCode: CountryGermany.Alpha2,
}
var martinUpdateContactInput = model.ContactInputXs{
	PhoneNumber:                  "017683464774",
	PhoneNumberCountryAlphaCode:  CountryGermany.Alpha2,
	PhoneNumber2:                 stringP("086823149540"),
	PhoneNumber2CountryAlphaCode: &CountryGermany.Alpha2,
	Web:                          stringP("www.web-martin.com"),
}

var albert = model.User{
	ID:         2,
	IdentityID: "ee8292f9-6e2c-4825-9d1e-ef6776e954b4",
	FirstName:  "Albert",
	LastName:   "Einstein",
	Gender:     model.GenderMan,
	Locale:     LocaleGermany,
	Address: &model.Address{
		Address:    "Albertstr. 23",
		City:       "Berlin",
		PostalCode: "10023",
		Country:    CountryGermany,
	},
}

var marco = model.User{
	ID:         3,
	IdentityID: "ee9292f9-6e2c-4825-9d1e-ef6776e954b4",
	FirstName:  "Marco",
	LastName:   "Italo",
	Gender:     model.GenderMan,
	Locale:     LocaleGermany,
	Address: &model.Address{
		Address:    "Leopoldstr. 4",
		City:       "München",
		PostalCode: "80804",
		Country:    CountryGermany,
	},
}

var sarah = model.User{
	ID:         4,
	IdentityID: "ee1292f9-6e2c-4825-9d1e-ef6776e954b4",
	FirstName:  "Sarah",
	LastName:   "Klein",
	Gender:     model.GenderWomen,
	Locale:     LocaleGermany,
	Address: &model.Address{
		Address:    "Schmitzstr. 16",
		City:       "München",
		PostalCode: "80804",
		Country:    CountryGermany,
	},
}

var elonLastName = "Howard"
var elonGender = model.GenderMan
var elonEmail = "dev.anonymous@vecomentman.de"
var elonAddressInput = model.AddressInput{
	Address:          "Badestr. 23",
	City:             "Genfo",
	PostalCode:       "83410",
	CountryAlphaCode: CountryGermany.Alpha2,
}
var elonContactInput = model.ContactInput{
	Email:                        elonEmail,
	PhoneNumber:                  "017683464774",
	PhoneNumberCountryAlphaCode:  CountryGermany.Alpha2,
	PhoneNumber2:                 nil,
	PhoneNumber2CountryAlphaCode: nil,
	Web:                          nil,
}

var elonAddress = model.Address{
	Address:    "Badestr. 23",
	City:       "Genfo",
	PostalCode: "83410",
	Country:    CountryGermany,
}
var elonContact = model.Contact{
	Email:               elonEmail,
	PhoneNumber:         "017683464774",
	PhoneNumberCountry:  CountryGermany,
	PhoneNumber2:        "",
	PhoneNumber2Country: nil,
	Web:                 "",
}

var elonCreateInput = model.UserCreateInput{
	Email:    elonEmail,
	Password: "RedUnicorn23!",
}

var elon = model.User{
	FirstName: "ElonNew",
	LastName:  elonLastName,
	Gender:    elonGender,
	Locale:    LocaleGermany,
	Address:   &elonAddress,
	Contact:   &elonContact,
}

var elonSso = sso.SsoClaims{
	Name:              "",
	GivenName:         "",
	FamilyName:        "",
	MiddleName:        "",
	Nickname:          "",
	PreferredUsername: "",
	Email:             elonEmail,
	EmailVerified:     false,
	UpdatedAt:         0,
	RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    elonEmail,
		Subject:   "ee3492f9-6e2c-4825-9d1e-ef6776e954b4",
		Audience:  nil,
		ExpiresAt: nil,
		NotBefore: nil,
		IssuedAt:  nil,
		ID:        "",
	},
}

var martinSso = sso.SsoClaims{
	Name:              "",
	GivenName:         "",
	FamilyName:        "",
	MiddleName:        "",
	Nickname:          "",
	PreferredUsername: "",
	Email:             martinEmail,
	EmailVerified:     true,
	UpdatedAt:         0,
	RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    martinEmail,
		Subject:   martin.IdentityID,
		Audience:  nil,
		ExpiresAt: nil,
		NotBefore: nil,
		IssuedAt:  nil,
		ID:        "",
	},
}

func NewUsers() User {
	user := User{
		Martin:           martin,
		Albert:           albert,
		Marco:            marco,
		MarcosPizzaUsers: nil,
		Sarah:            sarah,
	}

	user.MarcosPizzaUsers = append(user.MarcosPizzaUsers, user.Marco, user.Sarah)

	return user
}

func NewInputUsers() InputUser {
	userInput := InputUser{
		Elon: elonCreateInput,
	}

	return userInput
}

func NewUpdateInputUsers() UpdateInputUser {
	userUpdateInput := UpdateInputUser{
		Martin: martinUpdateInput,
	}

	return userUpdateInput
}

func NewCreationUsers() CreationUser {
	cachedUsers := CreationUser{
		Elon: elon,
	}

	return cachedUsers
}

func NewSsoUsers() SsoUser {
	ssoUsers := SsoUser{
		Elon:   elonSso,
		Martin: martinSso,
	}

	return ssoUsers
}

func (u *User) ResetMartin() {
	u.Martin = martin
}

func (u *User) ResetAlbert() {
	u.Albert = albert
}

func (u *User) ResetMarco() {
	u.Marco = marco
}

func (u *User) ResetMarcosPizzaUsers() {
	u.ResetMarco()
	u.ResetSarah()
	u.MarcosPizzaUsers = append(u.MarcosPizzaUsers, u.Marco, u.Sarah)
}

func (u *User) ResetSarah() {
	u.Sarah = sarah
}

func (u *InputUser) ResetElon() {
	u.Elon = elonCreateInput
}

func (u *CreationUser) ResetElon() {
	u.Elon = elon
}

// helper functions
func stringP(value string) *string {
	return &value
}
