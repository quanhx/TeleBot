package services

import (
	"xcheck.info/telebot/pkg/dtos"
	"xcheck.info/telebot/pkg/models"
	"xcheck.info/telebot/pkg/repositories"
)

type TransactionService interface {
	CreateTransaction(req dtos.TransactionRequest) error
	FindByID(id uint) *models.Transaction
	FindByUserID(userID uint) []models.Transaction
	GetBalance(userID uint) (float64, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}

}

func (s *transactionService) CreateTransaction(req dtos.TransactionRequest) error {
	var transaction = models.Transaction{}
	transaction = dtos.ToTransaction(req)
	return s.transactionRepo.CreateTransaction(&transaction)
}

func (s *transactionService) FindByID(id uint) *models.Transaction {
	var response = s.transactionRepo.FindByID(id)
	if response == nil {
		return nil
	}
	return response
}

func (s *transactionService) FindByUserID(userID uint) []models.Transaction {
	var response = s.transactionRepo.FindByUserID(userID)
	if response == nil {
		return nil
	}
	return response
}

func (s *transactionService) GetBalance(id uint) (float64, error) {
	var response, err = s.transactionRepo.GetBalance(id)
	if err != nil {
		return 0, err
	}
	return response, nil
}

func CreatePayment(req dtos.TransactionRequest) error {
	var transaction = models.Transaction{}
	var transactionRepo repositories.TransactionRepository
	transaction = dtos.ToTransaction(req)
	return transactionRepo.CreateTransaction(&transaction)
}