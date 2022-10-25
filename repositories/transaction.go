package repositories

import (
	"github.com/jinzhu/gorm"
	"xcheck.info/telebot/pkg/models"
)

type Balance struct {
	UserID  uint
	Balance float64
}

type TransactionRepository interface {
	FindByID(id uint) *models.Transaction
	FindByUserID(userID uint) []models.Transaction
	CreateTransaction(transaction *models.Transaction) error
	GetBalance(userID uint) (float64, error)
}

type transactionRepository struct {
	orm *gorm.DB
}

func NewTransactionRepository(orm *gorm.DB) TransactionRepository {
	return &transactionRepository{
		orm: orm,
	}
}

func (r *transactionRepository) FindByID(id uint) *models.Transaction {
	var transaction models.Transaction
	r.orm.First(&transaction, id)

	return &transaction
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.orm.Create(&transaction).Error
}

func (r *transactionRepository) FindByUserID(userID uint) []models.Transaction {
	var transactions []models.Transaction
	r.orm.Model(&models.Transaction{}).Where("status = ?", "finished").Where("deleted_at IS NULL").Where("user_id = ?", userID).Find(&transactions)

	return transactions
}

func (r *transactionRepository) GetBalance(userID uint) (float64, error) {
	var balance Balance
	//r.orm.Raw("select sum(case when payment_type = 0 and status = 'finished' then price_amount end)"+
	//	"- sum(case when payment_type = 1 and status = 'finished' then price_amount end)"+
	//	"from transactions where user_id = ?", userID).Scan(balance)

	//r.orm.Raw("select sum(price_amount) as balance from transactions where status = ?", "finished").Where("userID= ?", userID).Scan(balance)

	r.orm.Table("transactions").Select("sum(price_amount) as balance").Where("status =?", "finished").Where("user_id = ?", userID).Scan(&balance)
	return balance.Balance, nil
}
