package Common

import (
	"reflect"
)

const LORA = "lora"
const BLE = "ble"
const LPWA = "lpwa"
const ADDRESS = "address"
const DART = "DarTdi"

func FindFromJsonObj(json map[string]interface{}, key string) interface{} {
	var res interface{} = nil

	for k, v := range json {
		if k == key {
			return v
		}
		switch reflect.TypeOf(v).String() {
		case "map[string]interface {}":
			ret := FindFromJsonObj(v.(map[string]interface{}), key)
			if ret != nil {
				return ret
			}
			break
		default:
			break
		}
	}
	return res
}

func CheckFromJsonObj(json map[string]interface{}, key string, value interface{}) bool {
	var res = false

	for k, v := range json {
		switch reflect.TypeOf(v).String() {
		case "string":
			if k == key && v.(string) == value {
				return true
			}
			break
		case "float64":
			if k == key && v.(float64) == value {
				return true
			}
			break
		case "map[string]interface {}":
			ret := CheckFromJsonObj(v.(map[string]interface{}), key, value)
			if ret == true {
				return true
			}
			break
		default:
			break
		}
	}
	return res
}
