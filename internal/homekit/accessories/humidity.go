package accessories

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type SensorHumidity struct {
	*accessory.Accessory
	HumiditySensor *service.HumiditySensor
}

func NewSensorHumidity(info accessory.Info, _ ...interface{}) *SensorHumidity {
	acc := SensorHumidity{}
	acc.Accessory = accessory.New(info, accessory.TypeThermostat)
	acc.HumiditySensor = service.NewHumiditySensor()
	acc.AddService(acc.HumiditySensor.Service)
	return &acc
}

func (acc *SensorHumidity) GetType() uint8 {
	return uint8(acc.Accessory.Type)
}

func (acc *SensorHumidity) GetID() uint64 {
	return acc.Accessory.ID
}

func (acc *SensorHumidity) GetSN() string {
	return acc.Accessory.Info.SerialNumber.GetValue()
}

func (acc *SensorHumidity) GetName() string {
	return acc.Accessory.Info.Name.GetValue()
}

func (acc *SensorHumidity) GetAccessory() *accessory.Accessory {
	return acc.Accessory
}

func (acc *SensorHumidity) OnValuesRemoteUpdates(fn func()) {}
func (acc *SensorHumidity) OnExample()                      {}
