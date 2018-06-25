package db

import (
	"context"
	"database/sql"
	"github.com/adewaleafolabi/listing/model"
	"github.com/sirupsen/logrus"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) SaveProperty(ctx context.Context, property model.Property) error {
	_, err := r.db.Exec(`INSERT INTO property (id, title, body, amount) VALUES ($1,$2,$3,$4)`, property.ID, property.Title, property.Body, property.Amount)
	return err
}

func (r *PostgresRepository) ListProperties(ctx context.Context, limit int64, lastId string) ([]model.Property, error) {
	rows, err := r.db.Query(`SELECT id, title, sold, body, amount FROM property WHERE id > $1 LIMIT $2`, lastId, limit)
	if err != nil {
		return nil, err
	}
	var properties []model.Property

	for rows.Next() {
		property := model.Property{}
		if err = rows.Scan(&property.ID, &property.Title, &property.Sold, &property.Body, &property.Amount); err != nil {
			logrus.Error(err)
			fmt.Println(err)
			return nil, err
		}
		properties = append(properties, property)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *PostgresRepository) FindProperty(ctx context.Context, id string) (model.Property, error) {
	row:= r.db.QueryRow(`SELECT id, title, sold, body, amount FROM property WHERE id = $1 LIMIT 1`, id)

	var property model.Property

	if err:= row.Scan(&property.ID, &property.Title, &property.Sold, &property.Body, &property.Amount); err != nil {
		return property, err
	}
	return property, nil
}
