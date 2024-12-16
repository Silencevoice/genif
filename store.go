package store

import "context"

type Store[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	GetMultipleByID(ctx context.Context, ids []string) ([]*T, error)
	GetAll(ctx context.Context) ([]*T, error)
	Insert(ctx context.Context, id string, entity *T) (*T, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, entity *T) error
}
