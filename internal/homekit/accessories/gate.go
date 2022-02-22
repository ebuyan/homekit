package accessories

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Gate struct {
	*accessory.Accessory
	GarageDoorOpener *service.GarageDoorOpener
}

func NewGate(info accessory.Info, _ ...interface{}) *Gate {
	acc := Gate{}
	acc.Accessory = accessory.New(info, accessory.TypeGarageDoorOpener)
	acc.GarageDoorOpener = service.NewGarageDoorOpener()
	acc.AddService(acc.GarageDoorOpener.Service)
	return &acc
}

func (acc *Gate) GetType() uint8 {
	return uint8(acc.Accessory.Type)
}

func (acc *Gate) GetID() uint64 {
	return acc.Accessory.ID
}

func (acc *Gate) GetSN() string {
	return acc.Accessory.Info.SerialNumber.GetValue()
}

func (acc *Gate) GetName() string {
	return acc.Accessory.Info.Name.GetValue()
}

func (acc *Gate) GetAccessory() *accessory.Accessory {
	return acc.Accessory
}

func (acc *Gate) OnValuesRemoteUpdates(fn func()) {
	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(int) { fn() })
}

func (acc *Gate) OnExample() {
	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(func(v int) {})
}
