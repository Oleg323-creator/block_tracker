package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

// TO CONNECT IT WITH DB
type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SaveBlocksToDB(block int64) error {

	queryBuilder := squirrel.Update("blocks").
		Set("block_number", block).
		Where(squirrel.Eq{"id": 1})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Block number saved to db:", block)

	return nil
}
