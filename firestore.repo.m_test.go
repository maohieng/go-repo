package repo

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/maohieng/go-firebase"
	gofire "github.com/maohieng/go-firestore"
	"github.com/maohieng/go-price"
	"log"
	"os"
	"testing"
)

const TestMenuCollection = "test_menus"

type Menu struct {
	SimpleRepoEntity
	Name  string             `firestore:"name" json:"name"`
	Price price.MutablePrice `firestore:"price" json:"price"`
}

func (m Menu) TableName() string {
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

		crudrepo = NewFirestoreRepository(fStore, TestMenuCollection, log.Default()).(*FirestoreRepo)

		return m.Run()
	}())
}

func TestCreate(t *testing.T) {
	id, err := crudrepo.Create(ctx, &Menu{
		SimpleRepoEntity: SimpleRepoEntity{},
		Name:             "test 1",
		Price:            price.NewMutablePrice(price.NewFromInt(125, 100, "KHR")),
	})

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	toDeletedId = id

	t.Logf("Created: %s", id)
}

func TestUpdate(t *testing.T) {
	fv := make(map[string]any, 0)
	fv["price"] = price.NewMutablePrice(price.NewFromInt(3131, 100, "KHR"))
	err := crudrepo.Update(ctx, "WjnvkMFKfN7I4UbK5NQM", fv)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}

func TestGetOne(t *testing.T) {
	item := &Menu{}
	err := crudrepo.GetOne(ctx, "WjnvkMFKfN7I4UbK5NQM", item, true)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(item)
	t.Logf("GetOne: %s", string(itemJson))
}

func TestGetAll(t *testing.T) {
	results, err := crudrepo.GetAll(ctx, func() BaseRepoEntity {
		return &Menu{}
	}, false)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(results)
	t.Logf("GetAll: %s", string(itemJson))
}

func TestGetAllActive(t *testing.T) {
	results, err := crudrepo.GetAll(ctx, func() BaseRepoEntity {
		return &Menu{}
	}, true)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	itemJson, _ := json.Marshal(results)
	t.Logf("GetAll: %s", string(itemJson))
}

func TestCreateAll(t *testing.T) {
	menus := []*Menu{
		{
			SimpleRepoEntity: SimpleRepoEntity{},
			Name:             "all in 1",
		},
		{
			SimpleRepoEntity: SimpleRepoEntity{},
			Name:             "all in 2",
		},
	}

	entities := make([]BaseRepoEntity, 0, len(menus))
	for _, menu := range menus {
		entities = append(entities, menu)
	}

	ids, err := crudrepo.CreateAll(ctx, entities)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	t.Logf("CreateAll: %v", ids)
}

func recursivePages(t *testing.T, onlyActive bool, page gofire.Page, numb int) {
	resultPage, err := crudrepo.Paginate(ctx, page, func() BaseRepoEntity {
		return &Menu{}
	}, onlyActive)
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
	recursivePages(t, onlyActive, resultPage, numb)
}

func TestPaginateAll(t *testing.T) {
	recursivePages(t, false, gofire.Page{
		NextToken: "",
		Limit:     3,
	}, 1)
}

func TestPaginateActive(t *testing.T) {
	recursivePages(t, true, gofire.Page{
		NextToken: "",
		Limit:     3,
	}, 1)
}

func TestSoftDelete(t *testing.T) {
	err := crudrepo.SoftDelete(ctx, "WjnvkMFKfN7I4UbK5NQM")
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	// This delete depends on Create test func above
	err := crudrepo.Delete(ctx, toDeletedId)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
}
