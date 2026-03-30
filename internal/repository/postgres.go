package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"salary-calculator/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance.
// Implements CalculationRepository interface.
func NewPostgresRepository(db *sql.DB) CalculationRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(calc *models.Calculation) error {
	query := `
INSERT INTO calculations 
(gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, employer_total, mode, created_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW())
RETURNING id, created_at
`
	err := r.db.QueryRow(query,
		calc.GrossSalary,
		calc.NetSalary,
		calc.OPV,
		calc.IPN,
		calc.VOSMS,
		calc.SO,
		calc.SN,
		calc.OOSMS,
		calc.EmployerTotal,
		calc.Mode,
	).Scan(&calc.ID, &calc.CreatedAt)

	if err != nil {
		log.Printf("error saving calculation to database: %v", err)
		return fmt.Errorf("failed to save calculation: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetHistory(limit, offset int) ([]*models.Calculation, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query := `
SELECT id, gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, employer_total, mode, created_at
FROM calculations
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("error querying calculation history: %v", err)
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer rows.Close()

	var history []*models.Calculation
	for rows.Next() {
		var c models.Calculation
		if err := rows.Scan(
			&c.ID, &c.GrossSalary, &c.NetSalary,
			&c.OPV, &c.IPN, &c.VOSMS,
			&c.SO, &c.SN, &c.OOSMS,
			&c.EmployerTotal, &c.Mode, &c.CreatedAt,
		); err != nil {
			log.Printf("error scanning calculation row: %v", err)
			return nil, fmt.Errorf("failed to parse history: %w", err)
		}
		history = append(history, &c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error iterating calculation rows: %v", err)
		return nil, fmt.Errorf("failed to iterate history: %w", err)
	}

	return history, nil
}

func (r *PostgresRepository) GetByID(id int64) (*models.Calculation, error) {
	query := `
SELECT id, gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, employer_total, mode, created_at
FROM calculations
WHERE id = $1
`
	var c models.Calculation
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.GrossSalary, &c.NetSalary,
		&c.OPV, &c.IPN, &c.VOSMS,
		&c.SO, &c.SN, &c.OOSMS,
		&c.EmployerTotal, &c.Mode, &c.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		log.Printf("error retrieving calculation by id %d: %v", id, err)
		return nil, fmt.Errorf("failed to get calculation: %w", err)
	}

	return &c, nil
}
