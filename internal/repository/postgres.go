package repository

import (
	"database/sql"
	"errors"
	"salary-calculator/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(calc *models.Calculation) error {
	query := `
	INSERT INTO calculations 
	(gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, mode, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())
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
		calc.Mode,
	).Scan(&calc.ID, &calc.CreatedAt)

	return err
}

func (r *PostgresRepository) GetHistory(limit, offset int) ([]*models.Calculation, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `
	SELECT id, gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, mode, created_at
	FROM calculations
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.Calculation
	for rows.Next() {
		var c models.Calculation
		if err := rows.Scan(
			&c.ID, &c.GrossSalary, &c.NetSalary,
			&c.OPV, &c.IPN, &c.VOSMS,
			&c.SO, &c.SN, &c.OOSMS,
			&c.Mode, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, &c)
	}
	return history, nil
}

func (r *PostgresRepository) GetByID(id int64) (*models.Calculation, error) {
	query := `
	SELECT id, gross_salary, net_salary, opv, ipn, vosms, so, sn, oosms, mode, created_at
	FROM calculations
	WHERE id=$1
	`
	var c models.Calculation
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.GrossSalary, &c.NetSalary,
		&c.OPV, &c.IPN, &c.VOSMS,
		&c.SO, &c.SN, &c.OOSMS,
		&c.Mode, &c.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
