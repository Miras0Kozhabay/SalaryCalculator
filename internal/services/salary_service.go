package services

import (
	"errors"
	"salary-calculator/internal/calculator"
	"salary-calculator/internal/models"
	"salary-calculator/internal/repository"
)

type SalaryService struct {
	Calc       *calculator.Calculator
	Repository repository.CalculationRepository
}

func NewSalaryService(calc *calculator.Calculator, repo repository.CalculationRepository) *SalaryService {
	return &SalaryService{
		Calc:       calc,
		Repository: repo,
	}
}

// Calculate salary
func (s *SalaryService) Calculate(req *models.CalculateRequest) (*models.CalculateResponse, error) {

	if req.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	if req.Mode != "gross" && req.Mode != "net" {
		return nil, errors.New("mode must be 'gross' or 'net'")
	}

	var resp *models.CalculateResponse
	var err error

	if req.Mode == "gross" {
		resp, err = s.Calc.CalculateFromGross(req.Amount)
	} else {
		resp, err = s.Calc.CalculateFromNet(req.Amount)
	}

	if err != nil {
		return nil, err
	}

	// сохраняем в БД
	calcModel := &models.Calculation{
		GrossSalary: resp.GrossSalary,
		NetSalary:   resp.NetSalary,
		OPV:         resp.OPV,
		IPN:         resp.IPN,
		VOSMS:       resp.VOSMS,
		SO:          resp.SO,
		SN:          resp.SN,
		OOSMS:       resp.OOSMS,
		Mode:        req.Mode,
	}

	err = s.Repository.Save(calcModel)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get last calculations
func (s *SalaryService) GetHistory(limit, offset int) ([]*models.Calculation, error) {
	return s.Repository.GetHistory(limit, offset)
}
