package repo

import (
	"context"
)

type Page struct {
	Items     []BaseEntityType `json:"items"`
	NextToken string           `json:"nextToken"`
}

type CRUDRepository interface {
	Create(ctx context.Context, item BaseEntityType) (id string, err error)
	CreateAll(ctx context.Context, items []BaseEntityType) (ids []string, err error)
	Update(ctx context.Context, id string, fv map[string]interface{}) error
	GetOne(ctx context.Context, id string, item BaseEntityType) error
	GetAll(ctx context.Context, newItem func() BaseEntityType) ([]BaseEntityType, error)
	Paginate(ctx context.Context, limit int, startToken string, newItem func() BaseEntityType) (Page, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}
