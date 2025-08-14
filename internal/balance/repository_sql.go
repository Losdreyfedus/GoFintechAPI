package balance

import (
	"backend_path/internal/domain"
	"database/sql"
	"fmt"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) GetByUserID(userID int) (*domain.Balance, error) {
	query := `
		SELECT user_id, amount, last_updated_at
		FROM balances
		WHERE user_id = ?
	`

	balance := &domain.Balance{}
	err := r.db.QueryRow(query, userID).Scan(
		&balance.UserID,
		&balance.Amount,
		&balance.LastUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("balance not found")
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (r *sqlRepository) Update(balance *domain.Balance) error {
	// First try to update existing balance
	updateQuery := `
		UPDATE balances 
		SET amount = ?, last_updated_at = ?
		WHERE user_id = ?
	`

	result, err := r.db.Exec(updateQuery, balance.Amount, balance.LastUpdatedAt, balance.UserID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// If no rows were affected, insert new balance
	if rowsAffected == 0 {
		insertQuery := `
			INSERT INTO balances (user_id, amount, last_updated_at)
			VALUES (?, ?, ?)
		`

		_, err = r.db.Exec(insertQuery, balance.UserID, balance.Amount, balance.LastUpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert balance: %w", err)
		}
	}

	return nil
}
