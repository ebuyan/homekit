package openhab

type ItemType struct {
	ohType string
	tags   []string
}

func NewType(ohType string, tags []string) ItemType {
	return ItemType{
		ohType: ohType,
		tags:   tags,
	}
}

func (t ItemType) IsLight() bool {
	for _, tag := range t.tags {
		if tag == "Light" {
			return true
		}
	}
	return false
}

func (t ItemType) IsSwitch() bool {
	return t.ohType == "Switch"
}

func (t ItemType) IsGate() bool {
	for _, tag := range t.tags {
		if tag == "gate" {
			return true
		}
	}
	return false
}

func (t ItemType) IsTemp() bool {
	for _, tag := range t.tags {
		if tag == "temp" {
			return true
		}
	}
	return false
}

func (t ItemType) IsHumidity() bool {
	for _, tag := range t.tags {
		if tag == "hum" {
			return true
		}
	}
	return false
}
