package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"xcheck.info/telebot/pkg/dtos"
	"xcheck.info/telebot/pkg/services"
)

type UserAPI struct {
	UserService services.UserService
}

func ProvideUserAPI(p services.UserService) UserAPI {
	return UserAPI{UserService: p}
}

func (p *UserAPI) FindByID(c *gin.Context) {
	id, _ :=  strconv.Atoi(c.Param("id"))
	product := p.UserService.FindByID(uint(id))

	c.JSON(http.StatusOK, gin.H{"product": dtos.ToUserDTO(*product)})
}

func (p *UserAPI) Create(c *gin.Context) {
	var userDTO dtos.UserRequest
	err := c.BindJSON(&userDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	err = p.UserService.CreateUser(userDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": dtos.ToUser(userDTO)})
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


func FindByID() {
	url := "https://api.nowpayments.io/v1/min-amount?currency_from=nano&currency_to=usd"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", "4BC787P-WZHMJRJ-NBRA3K7-E31XHAV")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}