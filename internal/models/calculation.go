package models

import "time"

type Calculation struct {
	ID int64 `json:"id"`

	GrossSalary float64 `json:"gross_salary"`
	NetSalary   float64 `json:"net_salary"`

	OPV   float64 `json:"opv"`
	IPN   float64 `json:"ipn"`
	VOSMS float64 `json:"vosms"`

	SO    float64 `json:"so"`
	SN    float64 `json:"sn"`
	OOSMS float64 `json:"oosms"`

	Mode string `json:"mode"`

	CreatedAt time.Time `json:"created_at"`
}

type CalculateRequest struct {
	Amount float64 `json:"salary"`
	Mode   string  `json:"mode"`
}

type CalculateResponse struct {
	GrossSalary float64 `json:"gross_salary"`
	NetSalary   float64 `json:"net_salary"`

	OPV   float64 `json:"opv"`
	IPN   float64 `json:"ipn"`
	VOSMS float64 `json:"vosms"`

	SO    float64 `json:"so"`
	SN    float64 `json:"sn"`
	OOSMS float64 `json:"oosms"`
}
