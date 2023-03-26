package repo

const ActiveFieldName = "active"

type BaseEntityType interface {
	GetActive() bool
	SetActive(active bool)
	GetId() string
	SetId(id string)
	TableName() string
}

type BaseEntity struct {
	Active bool   `json:"-" firestore:"active" db:"active"`
	Id     string `json:"id" firestore:"-" db:"id"`
}

func (b *BaseEntity) GetActive() bool {
	return b.Active
}

func (b *BaseEntity) SetActive(active bool) {
	b.Active = active
}

func (b *BaseEntity) GetId() string {
	return b.Id
}

func (b *BaseEntity) SetId(id string) {
	b.Id = id
}
