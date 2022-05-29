package factory

import "gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"

type Company struct {
	MarcosPizza model.Company
	NewCompany  model.Company
}

var newCompany = model.Company{
	ID:   2,
	Name: "Johns Best Pizza",
	Address: &model.Address{
		Address:    "Leopoldstr. 49",
		City:       "Augsburg",
		PostalCode: "80333",
		Country:    &model.Country{ID: 65},
	},
	Contact: &model.Contact{
		Email:               "dev.john-pizza@vecomentman.de",
		PhoneNumber:         "08932455930",
		PhoneNumberCountry:  &model.Country{ID: 65},
		PhoneNumber2:        "08938889755",
		PhoneNumber2Country: &model.Country{ID: 65},
		Web:                 "https://johns-best-pizza.de",
	},
	TaxID: "DE987654321",
}

var marcosPizza = model.Company{
	ID:   1,
	Name: "Marcos Pizza",
	Address: &model.Address{
		Address:    "Brienner Str. 49",
		City:       "München",
		PostalCode: "80333",
		Country:    CountryGermany,
	},
	Contact: &model.Contact{
		Email:               "dev.marcos-pizza@vecomentman.de",
		PhoneNumber:         "08938889733",
		PhoneNumberCountry:  CountryGermany,
		PhoneNumber2:        "08938889755",
		PhoneNumber2Country: CountryGermany,
		Web:                 "https://marcospizza.de",
	},
	TaxID: "DE123456789",
}

func NewCompanies() Company {
	company := Company{
		MarcosPizza: marcosPizza,
		NewCompany:  newCompany,
	}

	return company
}

func (c *Company) ResetMarcosPizza() {
	c.MarcosPizza = marcosPizza
}

func (c *Company) ResetNewCompany() {
	c.NewCompany = newCompany
}
