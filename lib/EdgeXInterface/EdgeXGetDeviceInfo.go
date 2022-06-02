package EdgeXInterface

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface/URL"
)

func GetDeviceProtocol(deviceName string) map[string]interface{} {

	resp, err := http.Get(URL.GetDeviceInfoByName + "/" + deviceName)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	dev_info_json, _ := ioutil.ReadAll(resp.Body)
	dev_info := make(map[string]interface{})
	json.Unmarshal(dev_info_json, &dev_info)

	protocols := Common.FindFromJsonObj(dev_info, "protocols").(map[string]interface{})

	return protocols
	//res.BLE.Address = protocols["ble"].(map[string]interface{})["address"].(string)
	//res.LORA.Address = protocols["lora"].(map[string]interface{})["address"].(string)
}

func GetDeviceID(deviceName string) string {

	resp, err := http.Get(URL.GetDeviceInfoByName + "/" + deviceName)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	dev_info_json, _ := ioutil.ReadAll(resp.Body)
	dev_info := make(map[string]string)
	json.Unmarshal(dev_info_json, &dev_info)

	return dev_info["id"]
}

func GetDeviceName(deviceID string) string {
	resp, err := http.Get(URL.GetDeviceInfo + "/" + deviceID)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	dev_info_json, _ := ioutil.ReadAll(resp.Body)
	dev_info := make(map[string]string)
	json.Unmarshal(dev_info_json, &dev_info)

	return dev_info["name"]
}

func IsRedundantDevice(deviceAddress string) (result bool, previousDeviceID string) {

	resp, err := http.Get(URL.GetDeviceInfo)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	resp_str, _ := ioutil.ReadAll(resp.Body)
	var resp_json map[string]interface{}
	err = json.Unmarshal(resp_str, &resp_json)
	dev_info_list := resp_json["devices"].([]interface{})

	if err != nil {
		panic(err)
	}

	for _, v := range dev_info_list {
		if Common.CheckFromJsonObj(v.(map[string]interface{}), "address", deviceAddress) == true {
			//previousDeviceID = Common.FindFromJsonObj(v.(map[string]interface{}), "name").(string)
			previousDeviceID = v.(map[string]interface{})["name"].(string)
			return true, previousDeviceID
		}
	}

	return false, ""
}
