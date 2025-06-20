package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL database")

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS boxes (
		id UUID PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		items TEXT[],
		qr_code TEXT NOT NULL,
		qr_code_url TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_boxes_user_id ON boxes(user_id);
	CREATE INDEX IF NOT EXISTS idx_boxes_created_at ON boxes(created_at);
	
	-- Add description column if it doesn't exist (for existing databases)
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					   WHERE table_name='boxes' AND column_name='description') THEN
			ALTER TABLE boxes ADD COLUMN description TEXT;
		END IF;
	END $$;

	-- Add room column if it doesn't exist (for existing databases)
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					   WHERE table_name='boxes' AND column_name='room') THEN
			ALTER TABLE boxes ADD COLUMN room TEXT;
		END IF;
	END $$;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("✅ Database schema initialized successfully")
	return nil
}
