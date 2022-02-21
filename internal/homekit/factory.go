package homekit

import (
	"github.com/brutella/hc/accessory"
	openhabcli "github.com/ebuyan/ohyandex/pkg/openhab"
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
	f.checkStateAndUpdateSwitch(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildLight(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessory.NewLightbulb(info)
	ac.Lightbulb.On.SetValue(helper.StringToBool(item.State))
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) { f.repo.SetItemState(on, item) })
	f.checkStateAndUpdateLight(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildGate(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := NewAccessoryGate(info)
	ac.GarageDoorOpener.CurrentDoorState.SetValue(1)
	ac.GarageDoorOpener.TargetDoorState.SetValue(1)
	ac.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(on int) {
		ac.GarageDoorOpener.CurrentDoorState.SetValue(on)
		f.repo.SetItemState(true, item)
	})
	f.checkStateAndUpdateGate(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildTemp(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := accessory.NewTemperatureSensor(info, helper.GetFloatFromTemperature(item.State), 10.0, 30.0, 0.1)
	f.checkStateAndUpdateTemp(ac)
	return ac.Accessory
}

func (f AccessoryFactory) buildHum(item openhabcli.Item, id uint64) *accessory.Accessory {
	info := accessory.Info{Name: item.Label, ID: id}
	ac := NewAccessorySensorHumidity(info)
	ac.HumiditySensor.CurrentRelativeHumidity.SetValue(helper.GetFloatFromHumidity(item.State))
	f.checkStateAndUpdateHum(ac)
	return ac.Accessory
}

func (f AccessoryFactory) checkStateAndUpdateSwitch(ac *accessory.Switch) {
	go func() {
		for {
			val := helper.BoolToString(ac.Switch.On.GetValue())
			item, err := f.repo.GetItem(ac.Info.Name.GetValue())
			if err == nil && item.State != val {
				ac.Switch.On.SetValue(helper.StringToBool(item.State))
			}
			select {
			case <-time.After(1 * time.Second):
			}
		}
	}()
}

func (f AccessoryFactory) checkStateAndUpdateLight(ac *accessory.Lightbulb) {
	go func() {
		for {
			val := helper.BoolToString(ac.Lightbulb.On.GetValue())
			item, err := f.repo.GetItem(ac.Info.Name.GetValue())
			if err == nil && item.State != val {
				ac.Lightbulb.On.SetValue(helper.StringToBool(item.State))
			}
			select {
			case <-time.After(1 * time.Second):
			}
		}
	}()
}

func (f AccessoryFactory) checkStateAndUpdateGate(ac *AccessoryGate) {
	go func() {
		var gateTimeout = 0
		var maxTimeout = 60 * 2
		for {
			if item, err := f.repo.GetItem(ac.Info.Name.GetValue()); err == nil {
				if item.State == "ON" && gateTimeout == 0 {
					gateTimeout = 1
				}
				if gateTimeout != 0 {
					gateTimeout += 1
				}
				if gateTimeout > maxTimeout {
					ac.GarageDoorOpener.CurrentDoorState.SetValue(1)
					ac.GarageDoorOpener.TargetDoorState.SetValue(1)
					gateTimeout = 0
				}
			}
			select {
			case <-time.After(1 * time.Second):
			}
		}
	}()
}

func (f AccessoryFactory) checkStateAndUpdateTemp(ac *accessory.Thermometer) {
	go func() {
		for {
			if item, err := f.repo.GetItem(ac.Info.Name.GetValue()); err == nil {
				ac.TempSensor.CurrentTemperature.SetValue(helper.GetFloatFromTemperature(item.State))
			}
			select {
			case <-time.After(60 * 5 * time.Second):
			}
		}
	}()
}

func (f AccessoryFactory) checkStateAndUpdateHum(ac *AccessorySensorHumidity) {
	go func() {
		for {
			if item, err := f.repo.GetItem(ac.Info.Name.GetValue()); err == nil {
				ac.HumiditySensor.CurrentRelativeHumidity.SetValue(helper.GetFloatFromHumidity(item.State))
			}
			select {
			case <-time.After(60 * 5 * time.Second):
			}
		}
	}()
}
