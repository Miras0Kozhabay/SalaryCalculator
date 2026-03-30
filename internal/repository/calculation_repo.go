package repository

import "salary-calculator/internal/models"

// CalculationRepository defines the interface for calculation data access.
// This interface allows the service layer to depend on abstractions,
// not concrete implementations. This enables easy testing and future
// implementation changes.
type CalculationRepository interface {
	// Save stores a calculation record in the database.
	// Returns the calculation with ID and CreatedAt populated.
	Save(calc *models.Calculation) error

	// GetHistory retrieves recent calculation records with pagination.
	// limit: maximum number of results (defaults to 10 if <= 0)
	// offset: number of records to skip
	GetHistory(limit, offset int) ([]*models.Calculation, error)

	// GetByID retrieves a single calculation by ID.
	// Returns sql.ErrNoRows if not found.
	GetByID(id int64) (*models.Calculation, error)
}
