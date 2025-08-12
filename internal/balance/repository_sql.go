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
	query := `
		MERGE balances AS target
		USING (SELECT ? AS user_id, ? AS amount, ? AS last_updated_at) AS source
		ON target.user_id = source.user_id
		WHEN MATCHED THEN
			UPDATE SET amount = source.amount, last_updated_at = source.last_updated_at
		WHEN NOT MATCHED THEN
			INSERT (user_id, amount, last_updated_at) VALUES (source.user_id, source.amount, source.last_updated_at);
	`

	_, err := r.db.Exec(
		query,
		balance.UserID,
		balance.Amount,
		balance.LastUpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
