package transaction

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

func (r *sqlRepository) Create(tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		SELECT SCOPE_IDENTITY()
	`

	var id int
	err := r.db.QueryRow(
		query,
		tx.FromUserID,
		tx.ToUserID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.CreatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	tx.ID = id
	return nil
}

func (r *sqlRepository) GetByID(id int) (*domain.Transaction, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, type, status, created_at
		FROM transactions
		WHERE id = ?
	`

	tx := &domain.Transaction{}
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.FromUserID,
		&tx.ToUserID,
		&tx.Amount,
		&tx.Type,
		&tx.Status,
		&tx.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return tx, nil
}

func (r *sqlRepository) GetByUser(userID int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, type, status, created_at
		FROM transactions
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.FromUserID,
			&tx.ToUserID,
			&tx.Amount,
			&tx.Type,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *sqlRepository) Update(tx *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET from_user_id = ?, to_user_id = ?, amount = ?, type = ?, status = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		tx.FromUserID,
		tx.ToUserID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}
