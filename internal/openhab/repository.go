package openhab

import (
	"errors"
	openhabcli "github.com/ebuyan/ohyandex/pkg/openhab"
	"homekit/internal/pkg/helper"
	"log"
	"sync"
	"time"
)

type Repository struct {
	items  []openhabcli.Item
	config Config
	client openhabcli.Client
	sync.Mutex
}

func NewRepository(config Config, client openhabcli.Client) *Repository {
	repo := &Repository{
		config: config,
		client: client,
	}
	repo.loadItems()
	return repo
}

func (r *Repository) GetItems() []openhabcli.Item {
	for len(r.items) == 0 {
	}
	return r.items
}

func (r *Repository) GetItem(label string) (items openhabcli.Item, err error) {
	for _, item := range r.items {
		if item.Label == label {
			return item, nil
		}
	}
	err = errors.New("item not found " + label)
	return
}

func (r *Repository) SetItemState(val bool, item openhabcli.Item) {
	r.client.SetState(r.config.GetCredentials(), item.Name, helper.BoolToString(val))
	r.initItems()
}

func (r *Repository) loadItems() {
	go func() {
		for {
			prevCount := len(r.items)
			r.initItems()
			curCount := len(r.items)
			if prevCount != 0 && prevCount != curCount {
				log.Fatalln("new items. reloading app")
			}
			select {
			case <-time.After(10 * time.Second):
			}
		}
	}()
}

func (r *Repository) initItems() {
	r.Lock()
	r.items, _, _ = r.client.GetAllItemsByTag(r.config.GetCredentials(), r.config.searchTag)
	r.Unlock()
}
