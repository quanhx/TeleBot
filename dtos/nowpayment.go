package dtos

import "time"

type NowPayment struct {
	PaymentId        string    `json:"payment_id"`
	PaymentStatus    string    `json:"payment_status"`
	PayAddress       string    `json:"pay_address"`
	PriceAmount      float64   `json:"price_amount"`
	PriceCurrency    string    `json:"price_currency"`
	PayAmount        float64   `json:"pay_amount"`
	PayCurrency      string    `json:"pay_currency"`
	OrderId          string    `json:"order_id"`
	OrderDescription string    `json:"order_description"`
	PayInExtraId     string    `json:"payin_extra_id"`
	IpnCallbackUrl   string    `json:"ipn_callback_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PurchaseId       string    `json:"purchase_id"`
}
