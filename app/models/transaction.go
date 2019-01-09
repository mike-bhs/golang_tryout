package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model

	Amount   float64
	Currency string
}

type Transactions []*Transaction

func (Transaction) TableName() string {
	return "transactions"
}
