package db

import (
	"context"
	"github.com/adewaleafolabi/listing/model"
)

type Repository interface {
	SaveProperty(ctx context.Context, property model.Property) error
	ListProperties(ctx context.Context, limit int64, lastId string) ([]model.Property, error)
	FindProperty(ctx context.Context, id string) (model.Property, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}

func SaveProperty(ctx context.Context, property model.Property) error {
	return repository.SaveProperty(ctx, property)
}

func ListProperties(ctx context.Context, limit int64, lastId string) ([]model.Property, error) {
	return repository.ListProperties(ctx, limit, lastId)
}

func FindProperty(ctx context.Context, id string) (model.Property, error) {
	return repository.FindProperty(ctx, id)
}
