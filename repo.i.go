package repo

import (
	"context"
)

type CRUDRepository interface {
	Create(ctx context.Context, item BaseEntityType) (id string, err error)
	Update(ctx context.Context, id string, fv map[string]interface{}) error
	GetOne(ctx context.Context, id string, item BaseEntityType) error
	GetAll(ctx context.Context, newItem func() BaseEntityType) ([]BaseEntityType, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}
