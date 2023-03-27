package repo

import (
	"context"
	gofire "github.com/maohieng/go-firestore"
)

type CRUDRepository interface {
	Create(ctx context.Context, item BaseRepoEntity) (id string, err error)
	CreateAll(ctx context.Context, items []BaseRepoEntity) (ids []string, err error)
	Update(ctx context.Context, id string, fv map[string]interface{}) error
	GetOne(ctx context.Context, id string, item BaseRepoEntity, onlyActive bool) error
	GetAll(ctx context.Context, newItem func() BaseRepoEntity, onlyActive bool) ([]BaseRepoEntity, error)
	Paginate(ctx context.Context, prevp gofire.Page, newItem func() BaseRepoEntity, onlyActive bool) (gofire.Page, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

//type CRUDRepository2[T BaseRepoEntity] interface {
//	Create(ctx context.Context, item *T) (id string, err error)
//	CreateAll(ctx context.Context, items []*T) (ids []string, err error)
//	Update(ctx context.Context, id string, fv map[string]interface{}) error
//	GetOne(ctx context.Context, id string, item *T, onlyActive bool) error
//	GetAll(ctx context.Context, newItem func() *T, onlyActive bool) ([]*T, error)
//	Paginate(ctx context.Context, limit int, startToken string, newItem func() *T, onlyActive bool) (gofire.Page, error)
//	Delete(ctx context.Context, id string) error
//	SoftDelete(ctx context.Context, id string) error
//}
