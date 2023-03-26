package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/maohieng/go-firebase"
	"github.com/maohieng/go-repo"
	"log"
	"os"
	"testing"
)

const TestMenuCollection = "test_menus"

type Menu struct {
	repo.BaseEntity
	Name string `firestore:"name" json:"name"`
}

func (m *Menu) TableName() string {
	return TestMenuCollection
}

var (
	ctx         context.Context
	crudrepo    *FirestoreRepo
	toDeletedId string
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx = context.Background()
		fApp := firebase.InitAppDefault(ctx)
		fStore, err := fApp.Firestore(ctx)
		if err != nil {
			log.Println("cmd", "init", "firebase", "firestore", "err", err)
			os.Exit(1)
		}
		defer func(client *firestore.Client) {
			err := client.Close()
			if err != nil {
				log.Println("cmd", "exit", "closing", "firestore", "err", err)
			}
		}(fStore)

		crudrepo = &FirestoreRepo{client: fStore, cllName: TestMenuCollection, logger: log.Default()}

		return m.Run()
	}())
}

func TestCreate(t *testing.T) {
	id, err := crudrepo.Create(ctx, &Menu{
		BaseEntity: repo.BaseEntity{},
		Name:       "test 1",
	})

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	toDeletedId = id

	t.Logf("Created: %s", id)
}

func TestUpdate(t *testing.T) {
	fv := make(map[string]any, 0)
	fv["name"] = "Update test again"
	err := crudrepo.Update(ctx, "WjnvkMFKfN7I4UbK5NQM", fv)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}

func TestGetOne(t *testing.T) {
	item := &Menu{}
	err := crudrepo.GetOne(ctx, "WjnvkMFKfN7I4UbK5NQM", item)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(item)
	t.Logf("GetOne: %s", string(itemJson))
}

func TestGetAll(t *testing.T) {
	results, err := crudrepo.GetAll(ctx, func() repo.BaseEntityType {
		return &Menu{}
	})

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(results)
	t.Logf("GetAll: %s", string(itemJson))
}

func TestCreateAll(t *testing.T) {
	menus := []*Menu{
		{
			BaseEntity: repo.BaseEntity{},
			Name:       "all in 1",
		},
		{
			BaseEntity: repo.BaseEntity{},
			Name:       "all in 2",
		},
	}

	entities := make([]repo.BaseEntityType, 0, len(menus))
	for _, menu := range menus {
		entities = append(entities, menu)
	}

	ids, err := crudrepo.CreateAll(ctx, entities)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	t.Logf("CreateAll: %v", ids)
}

func recursivePages(t *testing.T, page repo.Page, numb int) {
	resultPage, err := crudrepo.Paginate(ctx, 3, page.NextToken, func() repo.BaseEntityType {
		return &Menu{}
	})
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(resultPage)
	t.Logf("Page %d: %s", numb, string(itemJson))

	if resultPage.NextToken == "" {
		t.Log("Done pagination")
		return
	}

	numb++
	recursivePages(t, resultPage, numb)
}

func TestPaginate(t *testing.T) {
	recursivePages(t, repo.Page{}, 1)
}

func TestSoftDelete(t *testing.T) {
	err := crudrepo.SoftDelete(ctx, "WjnvkMFKfN7I4UbK5NQM")
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	err := crudrepo.SoftDelete(ctx, toDeletedId)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}
