package transaction

import (
	"backend_path/internal/domain"
	"context"
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
		OUTPUT INSERTED.id
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// Handle system transactions (from_user_id or to_user_id = -1)
	var fromUserID, toUserID sql.NullInt64

	if tx.FromUserID == -1 {
		fromUserID.Valid = false
	} else {
		fromUserID.Int64 = int64(tx.FromUserID)
		fromUserID.Valid = true
	}

	if tx.ToUserID == -1 {
		toUserID.Valid = false
	} else {
		toUserID.Int64 = int64(tx.ToUserID)
		toUserID.Valid = true
	}

	var id int
	err := r.db.QueryRow(
		query,
		fromUserID,
		toUserID,
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

func (r *sqlRepository) CreateWithTransaction(ctx context.Context, tx *domain.Transaction) error {
	// Begin database transaction
	dbTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			dbTx.Rollback()
		}
	}()

	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at)
		OUTPUT INSERTED.id
		VALUES (?, ?, ?, ?, ?, ?)
	`

	var id int
	err = dbTx.QueryRowContext(ctx,
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

	// Commit transaction
	if err = dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *sqlRepository) GetByID(id int) (*domain.Transaction, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, type, status, created_at
		FROM transactions
		WHERE id = ?
	`

	tx := &domain.Transaction{}
	var fromUserID, toUserID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&fromUserID,
		&toUserID,
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

	// Handle NULL values
	if fromUserID.Valid {
		tx.FromUserID = int(fromUserID.Int64)
	} else {
		tx.FromUserID = -1 // System transaction
	}

	if toUserID.Valid {
		tx.ToUserID = int(toUserID.Int64)
	} else {
		tx.ToUserID = -1 // System transaction
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
		var fromUserID, toUserID sql.NullInt64

		err := rows.Scan(
			&tx.ID,
			&fromUserID,
			&toUserID,
			&tx.Amount,
			&tx.Type,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Handle NULL values
		if fromUserID.Valid {
			tx.FromUserID = int(fromUserID.Int64)
		} else {
			tx.FromUserID = -1 // System transaction
		}

		if toUserID.Valid {
			tx.ToUserID = int(toUserID.Int64)
		} else {
			tx.ToUserID = -1 // System transaction
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

	// Handle system transactions (from_user_id or to_user_id = -1)
	var fromUserID, toUserID sql.NullInt64

	if tx.FromUserID == -1 {
		fromUserID.Valid = false
	} else {
		fromUserID.Int64 = int64(tx.FromUserID)
		fromUserID.Valid = true
	}

	if tx.ToUserID == -1 {
		toUserID.Valid = false
	} else {
		toUserID.Int64 = int64(tx.ToUserID)
		toUserID.Valid = true
	}

	_, err := r.db.Exec(
		query,
		fromUserID,
		toUserID,
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
