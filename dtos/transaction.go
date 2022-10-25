package dtos

import (
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/models"
)

type TransactionRequest struct {
	gorm.Model
	UserID           int64   `json:"user_id"`
	Status           string  `json:"status"`
	PaymentType      int     `json:"payment_type"`
	PaymentAddress   string  `json:"payment_address"`
	PriceAmount      float64 `json:"price_amount"`
	PriceCurrency    string  `json:"price_currency"`
	PayAmount        float64 `json:"pay_amount"`
	PayCurrency      string  `json:"pay_currency"`
	OrderID          string  `json:"order_id"`
	OrderDescription string  `json:"order_description"`
}

func ToTransaction(transaction TransactionRequest) models.Transaction {
	return models.Transaction{
		UserID:           transaction.UserID,
		Status:           transaction.Status,
		PaymentType:      transaction.PaymentType,
		PayAddress:       transaction.PaymentAddress,
		PriceAmount:      transaction.PriceAmount,
		PriceCurrency:    transaction.PriceCurrency,
		PayAmount:        transaction.PayAmount,
		PayCurrency:      transaction.PayCurrency,
		OrderID:          transaction.OrderID,
		OrderDescription: transaction.OrderDescription,
	}
}

func ToTransactionDTO(transaction models.Transaction) TransactionRequest {
	return TransactionRequest{
		UserID:           transaction.UserID,
		Status:           transaction.Status,
		PaymentType:      transaction.PaymentType,
		PaymentAddress:   transaction.PayAddress,
		PriceAmount:      transaction.PriceAmount,
		PriceCurrency:    transaction.PriceCurrency,
		PayAmount:        transaction.PayAmount,
		PayCurrency:      transaction.PayCurrency,
		OrderID:          transaction.OrderID,
		OrderDescription: transaction.OrderDescription,
	}
}
