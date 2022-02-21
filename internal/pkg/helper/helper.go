package helper

import (
	"regexp"
	"strconv"
)

func GetFloatFromTemperature(temp string) float64 {
	re := regexp.MustCompile("[0-9.]+")
	res := re.FindAllString(temp, -1)
	if len(res) == 0 {
		return 0
	}
	floatRes, _ := strconv.ParseFloat(res[0], 64)
	return floatRes
}

func GetFloatFromHumidity(hum string) float64 {
	floatRes, _ := strconv.ParseFloat(hum, 64)
	return floatRes
}

func BoolToString(val bool) string {
	if val {
		return "ON"
	} else {
		return "OFF"
	}
}

func StringToBool(val string) bool {
	if val == "ON" {
		return true
	} else {
		return false
	}
}
