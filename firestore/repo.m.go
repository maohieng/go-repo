package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/maohieng/errs"
	gofire "github.com/maohieng/go-firestore"
	"github.com/maohieng/go-repo"
	"log"
)

func NewFirestoreRepository(client *firestore.Client, cllName string, logger *log.Logger) *FirestoreRepo {
	return &FirestoreRepo{client: client, cllName: cllName, logger: logger}
}

type FirestoreRepo struct {
	client  *firestore.Client
	cllName string
	logger  *log.Logger
}

func (f *FirestoreRepo) Create(ctx context.Context, item repo.BaseEntityType) (id string, err error) {
	const op errs.Op = "firestore.Create"
	id, err = gofire.Create(ctx, f.client, item)
	if err != nil {
		return "", errs.New(err, op)
	}

	return
}

func (f *FirestoreRepo) CreateAll(ctx context.Context, items []repo.BaseEntityType) (ids []string, err error) {
	const op errs.Op = "firestore.CreateAll"
	ids, err = gofire.BulkCreate(ctx, f.client, items)
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

func (f *FirestoreRepo) GetOne(ctx context.Context, id string, item repo.BaseEntityType) error {
	var op errs.Op = "firestore.GetOne"
	err := gofire.GetOne(ctx, f.client, id, item)
	if err != nil {
		return errs.New(err, op)
	}

	return nil
}

func (f *FirestoreRepo) GetAll(ctx context.Context, newItem func() repo.BaseEntityType) ([]repo.BaseEntityType, error) {
	var op errs.Op = "firestore.GetAll"
	results, err := gofire.GetAll(ctx, f.client, newItem)
	if err != nil {
		return nil, errs.New(err, op)
	}

	return results, nil
}

func (f *FirestoreRepo) Paginate(ctx context.Context, limit int, startToken string, newItem func() repo.BaseEntityType) (repo.Page, error) {
	var op errs.Op = "firestore.Pagination"
	p, err := gofire.Paginate(ctx, f.client, limit, startToken, newItem)
	if err != nil {
		return repo.Page{}, errs.New(err, op)
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
	err := gofire.SoftDelete(ctx, f.client, f.cllName, id)
	if err != nil {
		return errs.New(err, op)
	}

	return nil
}
