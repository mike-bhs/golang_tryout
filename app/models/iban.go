package models

import "time"

type Iban struct {
	UUID                   string    `json:"uuid"`
	AccountHolderName      string    `json:"account_holder_name"`
	AccountID              string    `json:"account_id"`
	BankInstitutionAddress string    `json:"bank_institution_address"`
	BankInstitutionCountry string    `json:"bank_institution_country"`
	BankInstitutionName    string    `json:"bank_institution_name"`
	BicSwift               string    `json:"bic_swift"`
	Currency               string    `json:"currency"`
	HouseAccountID         string    `json:"house_account_id"`
	IbanCode               string    `json:"iban_code"`
	IsEnabled              bool      `json:"enabled"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type Ibans []*Iban
