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
		INSERT INTO boxes (id, user_id, name, description, room, items, qr_code, qr_code_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		query,
		box.ID,
		box.UserID,
		box.Name,
		box.Description,
		box.Room,
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
		SELECT id, user_id, name, description, room, items, qr_code, qr_code_url, created_at, updated_at
		FROM boxes
		WHERE id = $1
	`

	box := &models.Box{}
	var items pq.StringArray
	var description sql.NullString
	var room sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&box.ID,
		&box.UserID,
		&box.Name,
		&description,
		&room,
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

	// Handle NULL description
	if description.Valid {
		box.Description = description.String
	} else {
		box.Description = ""
	}

	// Handle NULL room
	if room.Valid {
		box.Room = room.String
	} else {
		box.Room = ""
	}

	box.Items = []string(items)
	return box, nil
}

func (r *BoxRepository) GetByUserID(userID string) ([]*models.Box, error) {
	query := `
		SELECT id, user_id, name, description, room, items, qr_code, qr_code_url, created_at, updated_at
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
		var description sql.NullString
		var room sql.NullString

		err := rows.Scan(
			&box.ID,
			&box.UserID,
			&box.Name,
			&description,
			&room,
			&items,
			&box.QRCode,
			&box.QRCodeURL,
			&box.CreatedAt,
			&box.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan box: %w", err)
		}

		// Handle NULL description
		if description.Valid {
			box.Description = description.String
		} else {
			box.Description = ""
		}

		// Handle NULL room
		if room.Valid {
			box.Room = room.String
		} else {
			box.Room = ""
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
		SET name = $2, description = $3, room = $4, items = $5, updated_at = $6
		WHERE id = $1 AND user_id = $7
	`

	box.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		box.ID,
		box.Name,
		box.Description,
		box.Room,
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
