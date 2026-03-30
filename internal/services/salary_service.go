package services

import (
	"log"
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

// Calculate salary with validation and persistence
func (s *SalaryService) Calculate(req *models.CalculateRequest) (*models.CalculateResponse, error) {
	// Validate input
	if req.Amount <= 0 {
		return nil, ErrInvalidInput
	}

	if req.Mode != "gross" && req.Mode != "net" {
		return nil, ErrInvalidInput
	}

	// Perform calculation (may return ErrCalculation if algorithm fails)
	var resp *models.CalculateResponse
	var err error

	if req.Mode == "gross" {
		resp, err = s.Calc.CalculateFromGross(req.Amount)
	} else {
		resp, err = s.Calc.CalculateFromNet(req.Amount)
	}

	if err != nil {
		// Wrap calculation errors
		log.Printf("calculation error: %v", err)
		return nil, ErrCalculation
	}

	// Save to database
	calcModel := &models.Calculation{
		GrossSalary:   resp.GrossSalary,
		NetSalary:     resp.NetSalary,
		OPV:           resp.OPV,
		IPN:           resp.IPN,
		VOSMS:         resp.VOSMS,
		SO:            resp.SO,
		SN:            resp.SN,
		OOSMS:         resp.OOSMS,
		EmployerTotal: resp.EmployerTotal,
		Mode:          req.Mode,
	}

	if err := s.Repository.Save(calcModel); err != nil {
		log.Printf("database error: %v", err)
		return nil, ErrCalculation
	}

	return resp, nil
}

// Get last calculations
func (s *SalaryService) GetHistory(limit, offset int) ([]*models.Calculation, error) {
	return s.Repository.GetHistory(limit, offset)
}
