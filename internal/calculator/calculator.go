package calculator

import (
	"errors"
	"salary-calculator/internal/models"
)

type Calculator struct {
	MCI float64
}

// NewCalculator создает калькулятор с заданным МРП
func NewCalculator(mci float64) *Calculator {
	return &Calculator{MCI: mci}
}

// CalculateFromGross — из gross salary
func (c *Calculator) CalculateFromGross(gross float64) (*models.CalculateResponse, error) {
	if gross <= 0 {
		return nil, errors.New("gross salary must be > 0")
	}

	opv := gross * 0.10
	vosms := gross * 0.02
	ipnBase := gross - opv - vosms - 14*c.MCI
	ipn := ipnBase * 0.10
	net := gross - opv - vosms - ipn

	// работодателю
	so := (gross - opv) * 0.035
	oosms := gross * 0.03
	sn := (gross-opv-vosms)*0.095 - so
	if sn < 0 {
		sn = 0
	}

	resp := &models.CalculateResponse{
		GrossSalary: gross,
		NetSalary:   net,
		OPV:         opv,
		IPN:         ipn,
		VOSMS:       vosms,
		SO:          so,
		OOSMS:       oosms,
		SN:          sn,
	}

	return resp, nil
}

// CalculateFromNet — из net salary
func (c *Calculator) CalculateFromNet(net float64) (*models.CalculateResponse, error) {
	if net <= 0 {
		return nil, errors.New("net salary must be > 0")
	}

	// Используем итерацию, т.к. из net → gross сложнее
	gross := net
	var resp *models.CalculateResponse
	for i := 0; i < 10; i++ {
		var err error
		resp, err = c.CalculateFromGross(gross)
		if err != nil {
			return nil, err
		}

		diff := net - resp.NetSalary
		if diff < 0.01 && diff > -0.01 {
			break
		}
		gross += diff
	}

	return resp, nil
}
