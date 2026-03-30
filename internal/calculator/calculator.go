package calculator

import (
	"fmt"
	"log"
	"salary-calculator/internal/models"
)

const (
	MinSalary     = 90000
	MaxSalary     = 100000000
	MaxIterations = 20
	Epsilon       = 0.01 // 0.01 tenge precision
)

type Calculator struct {
	MCI float64
}

// NewCalculator creates a calculator with the specified MCI (МРП)
func NewCalculator(mci float64) *Calculator {
	return &Calculator{MCI: mci}
}

// CalculateFromGross calculates take-home salary from gross salary (before tax deductions)
func (c *Calculator) CalculateFromGross(gross float64) (*models.CalculateResponse, error) {
	if gross <= MinSalary {
		return nil, fmt.Errorf("salary is too small (minimum: %d, got: %.2f)", MinSalary, gross)
	}
	if gross > MaxSalary {
		return nil, fmt.Errorf("salary is too large (maximum: %d, got: %.2f)", MaxSalary, gross)
	}

	// Employee deductions
	opv := gross * 0.10
	vosms := gross * 0.02

	// IPN base = gross - OPV - VOSMS - 14*MCI
	ipnBase := gross - opv - vosms - 14*c.MCI
	if ipnBase < 0 {
		ipnBase = 0
	}
	ipn := ipnBase * 0.10

	// Net salary (take-home)
	net := gross - opv - vosms - ipn

	// Employer contributions
	so := (gross - opv) * 0.035
	oosms := gross * 0.03

	// SN = 9.5% of (gross - OPV - VOSMS) - SO
	snBase := (gross - opv - vosms) * 0.095
	sn := snBase - so
	if sn < 0 {
		sn = 0
	}

	resp := &models.CalculateResponse{
		GrossSalary:   gross,
		NetSalary:     net,
		OPV:           opv,
		IPN:           ipn,
		VOSMS:         vosms,
		SO:            so,
		OOSMS:         oosms,
		SN:            sn,
		EmployerTotal: gross + so + oosms + sn,
	}

	return resp, nil
}

// CalculateFromNet calculates gross salary from net (take-home) salary.
// Uses iterative approach with guaranteed convergence or error.
func (c *Calculator) CalculateFromNet(net float64) (*models.CalculateResponse, error) {
	if net <= MinSalary {
		return nil, fmt.Errorf("salary is too small (minimum: %d, got: %.2f)", MinSalary, net)
	}
	if net > MaxSalary {
		return nil, fmt.Errorf("salary is too large (maximum: %d, got: %.2f)", MaxSalary, net)
	}

	// Start with net as initial estimate for gross
	gross := net

	var resp *models.CalculateResponse
	for i := 0; i < MaxIterations; i++ {
		var err error
		resp, err = c.CalculateFromGross(gross)
		if err != nil {
			return nil, err
		}

		// Check if we've converged
		diff := net - resp.NetSalary
		if diff < Epsilon && diff > -Epsilon {
			log.Printf("CalculateFromNet converged after %d iterations (diff: %.2f)", i+1, diff)
			return resp, nil
		}

		// Adjust gross for next iteration
		gross += diff
	}

	// Did not converge - this is an error condition
	log.Printf("⚠️  Warning: CalculateFromNet did not converge after %d iterations", MaxIterations)
	return nil, fmt.Errorf(
		"calculation did not converge: net salary %.2f could not be reached (final diff: %.2f)",
		net, net-resp.NetSalary,
	)
}
