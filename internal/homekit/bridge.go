package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"homekit/internal/openhab"
	"log"
	"os"
)

type Bridge struct {
	openhabRepository *openhab.Repository
	factory           AccessoryFactory
	store             Store
}

func NewBridge(openhabRepository *openhab.Repository, factory AccessoryFactory, store Store) Bridge {
	return Bridge{
		openhabRepository: openhabRepository,
		factory:           factory,
		store:             store,
	}
}

func (b Bridge) Start() (err error) {
	acs, err := b.GetAccessories()
	if err != nil {
		return err
	}

	config := hc.Config{Pin: os.Getenv("PIN"), Port: os.Getenv("PORT")}
	bridge := accessory.NewBridge(accessory.Info{Name: "Openhab", ID: 1})
	t, err := hc.NewIPTransport(config, bridge.Accessory, acs...)
	if err != nil {
		return err
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
	return
}

func (b Bridge) GetAccessories() (acs []*accessory.Accessory, err error) {
	log.Println("loading items")
	items := b.openhabRepository.GetItems()
	if err != nil {
		return nil, err
	}
	for key, item := range items {
		ok, id := b.store.GetId(item.Name)
		if !ok {
			id = uint64(key) + 2
			b.store.SetId(id, item.Name)
		}
		log.Println("building item " + item.Name + " type " + item.Type)
		if ac := b.factory.Build(item, id); ac != nil {
			acs = append(acs, ac)
		}
	}
	return
}
