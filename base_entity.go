package repo

import gofire "github.com/maohieng/go-firestore"

const ActiveFieldName = "active"

type BaseRepoEntity interface {
	gofire.BaseEntity
	GetActive() bool
	SetActive(active bool)
}

type SimpleRepoEntity struct {
	Active bool   `json:"-" firestore:"active" db:"active"`
	Id     string `json:"id" firestore:"-" db:"id"`
}

func (b *SimpleRepoEntity) GetActive() bool {
	return b.Active
}

func (b *SimpleRepoEntity) SetActive(active bool) {
	b.Active = active
}

func (b *SimpleRepoEntity) GetId() string {
	return b.Id
}

func (b *SimpleRepoEntity) SetId(id string) {
	b.Id = id
}
