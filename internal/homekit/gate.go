package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type AccessoryGate struct {
	*accessory.Accessory
	GarageDoorOpener *service.GarageDoorOpener
}

func NewAccessoryGate(info accessory.Info, _ ...interface{}) *AccessoryGate {
	acc := AccessoryGate{}
	acc.Accessory = accessory.New(info, accessory.TypeGarageDoorOpener)
	acc.GarageDoorOpener = service.NewGarageDoorOpener()
	acc.AddService(acc.GarageDoorOpener.Service)
	return &acc
}

func (acc *AccessoryGate) GetType() uint8 {
	return uint8(acc.Accessory.Type)
}

func (acc *AccessoryGate) GetID() uint64 {
	return acc.Accessory.ID
}

func (acc *AccessoryGate) GetSN() string {
	return acc.Accessory.Info.SerialNumber.GetValue()
}

func (acc *AccessoryGate) GetName() string {
	return acc.Accessory.Info.Name.GetValue()
}

func (acc *AccessoryGate) GetAccessory() *accessory.Accessory {
	return acc.Accessory
}

func (acc *AccessoryGate) OnValuesRemoteUpdates(fn func()) {
	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(int) { fn() })
}

func (acc *AccessoryGate) OnExample() {
	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(v int) {})
}
