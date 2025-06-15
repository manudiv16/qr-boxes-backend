package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/qr-boxes/backend/internal/database"
	"github.com/qr-boxes/backend/internal/models"
)

type BoxRepository struct {
	db *database.DB
}

func NewBoxRepository(db *database.DB) *BoxRepository {
	return &BoxRepository{db: db}
}

func (r *BoxRepository) Create(box *models.Box) error {
	query := `
		INSERT INTO boxes (id, user_id, name, items, qr_code, qr_code_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Exec(
		query,
		box.ID,
		box.UserID,
		box.Name,
		pq.Array(box.Items),
		box.QRCode,
		box.QRCodeURL,
		box.CreatedAt,
		box.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create box: %w", err)
	}
	
	return nil
}

func (r *BoxRepository) GetByID(id string) (*models.Box, error) {
	query := `
		SELECT id, user_id, name, items, qr_code, qr_code_url, created_at, updated_at
		FROM boxes
		WHERE id = $1
	`
	
	box := &models.Box{}
	var items pq.StringArray
	
	err := r.db.QueryRow(query, id).Scan(
		&box.ID,
		&box.UserID,
		&box.Name,
		&items,
		&box.QRCode,
		&box.QRCodeURL,
		&box.CreatedAt,
		&box.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("box not found")
		}
		return nil, fmt.Errorf("failed to get box: %w", err)
	}
	
	box.Items = []string(items)
	return box, nil
}

func (r *BoxRepository) GetByUserID(userID string) ([]*models.Box, error) {
	query := `
		SELECT id, user_id, name, items, qr_code, qr_code_url, created_at, updated_at
		FROM boxes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user boxes: %w", err)
	}
	defer rows.Close()
	
	var boxes []*models.Box
	
	for rows.Next() {
		box := &models.Box{}
		var items pq.StringArray
		
		err := rows.Scan(
			&box.ID,
			&box.UserID,
			&box.Name,
			&items,
			&box.QRCode,
			&box.QRCodeURL,
			&box.CreatedAt,
			&box.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan box: %w", err)
		}
		
		box.Items = []string(items)
		boxes = append(boxes, box)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return boxes, nil
}

func (r *BoxRepository) Update(box *models.Box) error {
	query := `
		UPDATE boxes
		SET name = $2, items = $3, updated_at = $4
		WHERE id = $1 AND user_id = $5
	`
	
	box.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		box.ID,
		box.Name,
		pq.Array(box.Items),
		box.UpdatedAt,
		box.UserID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update box: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("box not found or unauthorized")
	}
	
	return nil
}

func (r *BoxRepository) Delete(id, userID string) error {
	query := `DELETE FROM boxes WHERE id = $1 AND user_id = $2`
	
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete box: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("box not found or unauthorized")
	}
	
	return nil
}

func (r *BoxRepository) GetUserBoxCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM boxes WHERE user_id = $1`
	
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get user box count: %w", err)
	}
	
	return count, nil
}