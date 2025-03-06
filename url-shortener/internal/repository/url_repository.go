package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener/internal/models"
)

type urlRepository struct {
	db *pgxpool.Pool
}

func NewUrlRepository(db *pgxpool.Pool) *urlRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO urls (url, address) VALUES ($1, $2)",
		url.URL,
		url.Address)
	if err != nil {
		return err
	}
	return nil
}

func (r *urlRepository) GetByURL(url string) (*models.URL, error) {
	var model models.URL
	row := r.db.QueryRow(context.Background(), "SELECT id, url, address FROM urls WHERE url = $1", url)
	if err := row.Scan(&model.Id, &model.URL, &model.Address); err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *urlRepository) GetByAddress(address string) (*models.URL, error) {
	var model models.URL
	row := r.db.QueryRow(context.Background(), "SELECT id, url, address FROM urls WHERE address = $1", address)
	if err := row.Scan(&model.Id, &model.URL, &model.Address); err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *urlRepository) DeleteByURL(url string) error {
	tag, err := r.db.Exec(context.Background(), "DELETE FROM urls WHERE url = $1", url)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
