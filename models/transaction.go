package models

import "github.com/jinzhu/gorm"

type Transaction struct {
	gorm.Model
	UserID           int64   `json:"user_id"`
	Status           string  `json:"status"`
	PaymentType      int     `json:"payment_type"`
	PayAddress       string  `json:"pay_address"`
	PriceAmount      float64 `json:"price_amount"`
	PriceCurrency    string  `json:"price_currency"`
	PayAmount        float64 `json:"pay_amount"`
	PayCurrency      string  `json:"pay_currency"`
	OrderID          string  `json:"order_id"`
	OrderDescription string  `json:"order_description"`
}
