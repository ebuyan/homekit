package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	openhabcli "github.com/ebuyan/ohyandex/pkg/openhab"
	"homekit/internal/homekit/accessories"
	"homekit/internal/openhab"
	"homekit/internal/pkg/helper"
	"log"
	"time"
)

type AccessoryFactory struct {
	repo *openhab.Repository
}

func NewFactory(repo *openhab.Repository) AccessoryFactory {
	return AccessoryFactory{
		repo: repo,
	}
}

func (f AccessoryFactory) Build(item openhabcli.Item, id uint64) *accessory.Accessory {
	itemType := openhab.NewType(item.Type, item.Tags)
	if itemType.IsGate() {
		return f.buildGate(item, id)
	}
	if itemType.IsLight() {
		return f.buildLight(item, id)
	}
	if itemType.IsSwitch() {
		return f.buildSwitch(item, id)
	}
	if itemType.IsHumidity() {
		return f.buildHum(item, id)
	}
	if itemType.IsTemp() {
		return f.buildTemp(item, id)
	}
	log.Println("unsupported type " + item.Type)
	return nil
}

func (f AccessoryFactory) buildSwitch(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessory.NewSwitch(info)
	ac.Switch.On.SetValue(helper.StringToBool(item.State))
	ac.Switch.On.OnValueRemoteUpdate(func(on bool) { f.repo.SetItemState(on, item) })
	go f.checkStateAndUpdateSwitch(ac.Switch.On, ac.Info.Name.GetValue())
	return ac.Accessory
}

func (f AccessoryFactory) buildLight(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessory.NewLightbulb(info)
	ac.Lightbulb.On.SetValue(helper.StringToBool(item.State))
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) { f.repo.SetItemState(on, item) })
	go f.checkStateAndUpdateSwitch(ac.Lightbulb.On, ac.Info.Name.GetValue())
	return ac.Accessory
}

func (f AccessoryFactory) buildGate(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessories.NewGate(info)
	ac.GarageDoorOpener.CurrentDoorState.SetValue(1)
	ac.GarageDoorOpener.TargetDoorState.SetValue(1)
	ac.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(on int) {
		ac.GarageDoorOpener.CurrentDoorState.SetValue(on)
		f.repo.SetItemState(true, item)
	})
	go f.checkStateAndUpdateGate(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildTemp(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessory.NewTemperatureSensor(info, helper.GetFloatFromTemperature(item.State), 10.0, 30.0, 0.1)
	go f.checkStateAndUpdateTemp(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildHum(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessories.NewSensorHumidity(info)
	ac.HumiditySensor.CurrentRelativeHumidity.SetValue(helper.GetFloatFromHumidity(item.State))
	go f.checkStateAndUpdateHum(ac)
	return ac.Accessory
}

func (f AccessoryFactory) checkStateAndUpdateSwitch(on *characteristic.On, itemName string) {
	for {
		item, err := f.repo.GetItem(itemName)
		if err == nil && item.State != helper.BoolToString(on.GetValue()) {
			on.SetValue(helper.StringToBool(item.State))
		}
		select {
		case <-time.After(1 * time.Second):
		}
	}
}

func (f AccessoryFactory) checkStateAndUpdateGate(ac *accessories.Gate) {
	for {
		//пока нет состояния ворот, просто сетим "closed"
		ac.GarageDoorOpener.CurrentDoorState.SetValue(1)
		ac.GarageDoorOpener.TargetDoorState.SetValue(1)
		select {
		case <-time.After(20 * time.Second):
		}
	}
}

func (f AccessoryFactory) checkStateAndUpdateTemp(ac *accessory.Thermometer) {
	for {
		if item, err := f.repo.GetItem(ac.Info.Name.GetValue()); err == nil {
			ac.TempSensor.CurrentTemperature.SetValue(helper.GetFloatFromTemperature(item.State))
		}
		select {
		case <-time.After(60 * 5 * time.Second):
		}
	}
}

func (f AccessoryFactory) checkStateAndUpdateHum(ac *accessories.SensorHumidity) {
	for {
		if item, err := f.repo.GetItem(ac.Info.Name.GetValue()); err == nil {
			ac.HumiditySensor.CurrentRelativeHumidity.SetValue(helper.GetFloatFromHumidity(item.State))
		}
		select {
		case <-time.After(60 * 5 * time.Second):
		}
	}
}
