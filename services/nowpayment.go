package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const APIKey = "3AR34P9-A1CMS53-KQVWXNN-TPQ4VW5"

func Status() {
	url := "https://api.sandbox.nowpayments.io/v1/status"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
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

func Currencies() {
	url := "https://api.sandbox.nowpayments.io/v1/currencies"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", APIKey)

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

func PriceOfService() {
	url := "https://api.sandbox.nowpayments.io/v1/estimate?amount=5&currency_from=usd&currency_to=nano"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", APIKey)

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

func Payment(amount int, currency string) *PaymentResponse {
	url := "https://api.sandbox.nowpayments.io/v1/payment"
	method := "POST"
	orderID := strconv.FormatInt(time.Now().Unix(), 10)

	payloadFormat := fmt.Sprintf(`{
  "price_amount": %v,
  "price_currency": "usd",
  "pay_currency": "%v",
  "ipn_callback_url": "https://nowpayments.io",
  "order_id": "%v",
  "order_description": "Payment for check information service x 1",
  "case" : "partially_paid"
}`, amount, currency, orderID)
	payload := strings.NewReader(payloadFormat)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("x-api-key", "3AR34P9-A1CMS53-KQVWXNN-TPQ4VW5")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var result PaymentResponse
	_ = json.Unmarshal(body, &result)
	return &result
}

func PaymentStatus(paymentID string) *PaymentStatusResponse {
	url := fmt.Sprintf("https://api.sandbox.nowpayments.io/v1/payment/%v", paymentID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("x-api-key", APIKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var result PaymentStatusResponse
	_ = json.Unmarshal(body, &result)
	return &result
}

func GetMinPaymentAmount() {
	url := "https://api.sandbox.nowpayments.io/v1/min-amount?currency_from=nano&currency_to=usd"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", APIKey)

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

func GetEstimatedPrice(amount int, cFrom, cTo string) {
	url := fmt.Sprintf("https://api.sandbox.nowpayments.io/v1/estimate?amount=%v&currency_from=%s&currency_to=%s", amount, cFrom, cTo)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", APIKey)

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

type PaymentResponse struct {
	PaymentId        string    `json:"payment_id"`
	PaymentStatus    string    `json:"payment_status"`
	PayAddress       string    `json:"pay_address"`
	PriceAmount      int       `json:"price_amount"`
	PriceCurrency    string    `json:"price_currency"`
	PayAmount        string    `json:"pay_amount"`
	AmountReceived   float64   `json:"amount_received"`
	PayCurrency      string    `json:"pay_currency"`
	OrderId          string    `json:"order_id"`
	OrderDescription string    `json:"order_description"`
	IpnCallbackUrl   string    `json:"ipn_callback_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PurchaseId       string    `json:"purchase_id"`
	SmartContract    string    `json:"smart_contract"`
	Network          string    `json:"network"`
	NetworkPrecision string    `json:"network_precision"`
	TimeLimit        time.Time `json:"time_limit"`
}

type PaymentStatusResponse struct {
	PaymentId        int64     `json:"payment_id"`
	PaymentStatus    string    `json:"payment_status"`
	PayAddress       string    `json:"pay_address"`
	PriceAmount      float64   `json:"price_amount"`
	PriceCurrency    string    `json:"price_currency"`
	PayAmount        float64   `json:"pay_amount"`
	ActuallyPaid     float64   `json:"actually_paid"`
	PayCurrency      string    `json:"pay_currency"`
	OrderId          string    `json:"order_id"`
	OrderDescription string    `json:"order_description"`
	PurchaseId       string    `json:"purchase_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	OutcomeAmount    float64   `json:"outcome_amount"`
	OutcomeCurrency  string    `json:"outcome_currency"`
}
