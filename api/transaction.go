package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"xcheck.info/telebot/pkg/dtos"
	"xcheck.info/telebot/pkg/services"
)

type TransactionAPI struct {
	TransactionService services.TransactionService
}

func ProvideTransactionAPI(transaction services.TransactionService) TransactionAPI {
	return TransactionAPI{TransactionService: transaction}
}

func (t *TransactionAPI) FindByID(c *gin.Context) {
	id, _ :=  strconv.Atoi(c.Param("id"))
	product := t.TransactionService.FindByID(uint(id))

	c.JSON(http.StatusOK, gin.H{"product": dtos.ToTransactionDTO(*product)})
}

func (t *TransactionAPI) Create(c *gin.Context) {
	var transactionDTO dtos.TransactionRequest
	err := c.BindJSON(&transactionDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	err = t.TransactionService.CreateTransaction(transactionDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": dtos.ToTransaction(transactionDTO)})
}

//func (p *UserAPI) Update(c *gin.Context) {
//	var productDTO ProductDTO
//	err := c.BindJSON(&productDTO)
//	if err != nil {
//		log.Fatalln(err)
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	id, _ :=  strconv.Atoi(c.Param("id"))
//	product := p.UserService.FindByID(uint(id))
//	if product == (Product{}) {
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	product.Code = productDTO.Code
//	product.Price = productDTO.Price
//	p.UserService.Save(product)
//
//	c.Status(http.StatusOK)
//}
//
//func (p *UserAPI) Delete(c *gin.Context) {
//	id, _ :=  strconv.Atoi(c.Param("id"))
//	product := p.UserService.FindByID(uint(id))
//	if product == (Product{}) {
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	p.UserService.Delete(product)
//
//	c.Status(http.StatusOK)
//}
