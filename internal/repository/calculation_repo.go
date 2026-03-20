package repository

import "salary-calculator/internal/models"

type CalculationRepository interface {
	Save(calc *models.Calculation) error
	GetHistory(limit, offset int) ([]*models.Calculation, error)
	GetByID(id int64) (*models.Calculation, error)
}
