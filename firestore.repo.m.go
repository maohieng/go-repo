package repo

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/maohieng/errs"
	gofire "github.com/maohieng/go-firestore"
	"log"
)

// NewFirestoreRepository creates a *FirestoreRepo instance.
func NewFirestoreRepository(client *firestore.Client, cllName string, logger *log.Logger) CRUDRepository {
	// >Note we return CRUDRepository instead of FirestoreRepo to fully support repo package
	// >and be able to see error if API changed
	return &FirestoreRepo{client: client, cllName: cllName, logger: logger}
}

type FirestoreRepo struct {
	client  *firestore.Client
	cllName string
	logger  *log.Logger
}

func (f *FirestoreRepo) Create(ctx context.Context, item BaseRepoEntity) (id string, err error) {
	const op errs.Op = "firestore.Create"
	item.SetActive(true)
	id, err = gofire.Create(ctx, f.client, item)
	if err != nil {
		return "", errs.New(err, op)
	}

	return
}

func (f *FirestoreRepo) CreateAll(ctx context.Context, items []BaseRepoEntity) (ids []string, err error) {
	const op errs.Op = "firestore.CreateAll"

	ents := make([]gofire.BaseEntity, 0, len(items))
	for _, item := range items {
		item.SetActive(true)
		ent, _ := item.(gofire.BaseEntity)
		ents = append(ents, ent)
	}

	ids, err = gofire.BulkCreate(ctx, f.client, true, ents...)
	if ids == nil && err != nil {
		return nil, errs.New(err, op)
	}

	if f.logger != nil {
		f.logger.Println(op, err)
	}

	return ids, nil
}

func (f *FirestoreRepo) Update(ctx context.Context, id string, fv map[string]interface{}) error {
	const op errs.Op = "firestore.Update"
	err := gofire.Update(ctx, f.client, f.cllName, id, fv)
	if err != nil {
		return errs.New(err, op)
	}

	return nil
}

func (f *FirestoreRepo) GetOne(ctx context.Context, id string, item BaseRepoEntity, onlyActive bool) error {
	var op errs.Op = "firestore.GetOne"
	err := gofire.GetOne(ctx, f.client, id, item)
	if err != nil {
		return errs.New(err, op)
	}

	if onlyActive && !item.GetActive() {
		return gofire.ErrNotFound
	}

	return nil
}

func (f *FirestoreRepo) GetAll(ctx context.Context, newItem func() BaseRepoEntity, onlyActive bool) ([]BaseRepoEntity, error) {
	var op errs.Op = "firestore.GetAll"
	var ents []gofire.BaseEntity
	var err error
	if !onlyActive {
		ents, err = gofire.GetAll(ctx, f.client, func() gofire.BaseEntity {
			return newItem()
		})
	} else {
		ents, err = gofire.GetAll(ctx, f.client, func() gofire.BaseEntity {
			return newItem()
		}, gofire.Where{
			Path:  ActiveFieldName,
			Op:    "==",
			Value: true,
		})
	}

	if err != nil {
		return nil, errs.New(err, op)
	}

	results := make([]BaseRepoEntity, 0, len(ents))
	for _, ent := range ents {
		r, _ := ent.(BaseRepoEntity)
		results = append(results, r)
	}

	return results, nil
}

func (f *FirestoreRepo) Paginate(ctx context.Context, prevp gofire.Page, newItem func() BaseRepoEntity, onlyActive bool) (gofire.Page, error) {
	var op errs.Op = "firestore.Paginate"
	var p gofire.Page
	var err error
	if !onlyActive {
		p, err = gofire.Paginate(ctx, f.client, prevp, func() gofire.BaseEntity {
			return newItem()
		})
	} else {
		p, err = gofire.Paginate(ctx, f.client, prevp, func() gofire.BaseEntity {
			return newItem()
		}, gofire.Where{
			Path:  ActiveFieldName,
			Op:    "==",
			Value: true,
		})
	}

	if err != nil {
		return gofire.Page{}, errs.New(err, op)
	}

	return p, nil
}

func (f *FirestoreRepo) Delete(ctx context.Context, id string) error {
	var op errs.Op = "firestore.Delete"
	err := gofire.Delete(ctx, f.client, f.cllName, id)
	if err != nil {
		return errs.New(err, op)
	}

	return nil
}

func (f *FirestoreRepo) SoftDelete(ctx context.Context, id string) error {
	var op errs.Op = "firestore.SoftDelete"
	err := gofire.Update(ctx, f.client, f.cllName, id, map[string]any{"active": false})
	if err != nil {
		return errs.New(err, op)
	}

	return nil
}
