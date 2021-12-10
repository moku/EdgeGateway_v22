package driver

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface"
)

func MisMsgHandler(client mqtt.Client, mqtt_msg mqtt.Message) {
	log.Println("==================================================[Network] Microservice MQTT Handler")
	//Scheduler("", client, mqtt_msg)
}

func encode(data string) string { //downlink
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	m := make(map[string]interface{})
	m["confirmed"] = false
	m["fPort"] = 2
	m["data"] = sEnc
	m["timing"] = "IMMEDIATELY"
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	downlink_data := string(jsonMsg)
	return downlink_data
}

func decode(data []byte) []byte { //uplink
	var message map[string]interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		panic(err)
	}
	//devEUI := message["devEUI"].(string)
	msgData := message["data"].(string)
	decoded, _ := b64.StdEncoding.DecodeString(msgData)
	//uplink_data := string(decoded)
	return decoded
}

func MsgHandler(client mqtt.Client, mqtt_msg mqtt.Message) {
	log.Println("==================================================[Network] Device MQTT Handler")

	data := mqtt_msg.Payload()
	topic := mqtt_msg.Topic()

	if topic == "pipeLineDown" {
		//InitNode(data)
		//Controller()

	} else if strings.Contains(topic, "application") {
		// LoRa
		str_arr := strings.Split(topic, "/")
		msg_type := str_arr[4]

		if msg_type == "join" {
			log.Println("==================================================[Network] LoRa Device JOIN, DEV_EUI: " + str_arr[3])
		} else if msg_type == "rx" {
			decoded_data := decode(data)
			if !bytes.Contains(decoded_data, []byte("rsc=")) {
				ParseRequestMsg(decoded_data, client, Common.LORA, str_arr[3])
			} else {
				ParseResponseMsg(decoded_data)
			}
		}
	} else if topic == Common.BLE {
		// BLE
		if !bytes.Contains(data, []byte("rsc=")) {
			ParseRequestMsg(data, client, Common.BLE, "NA")
		} else {
			ParseResponseMsg(data)
		}
	}

}

func SendDeviceInfo(result string) {

	MqttClient.Publish("edgex", 0, false, result)
}

func SendDownlinkMsg(msg string, di string) {
	topic := ""

	protocols := EdgeXInterface.GetDeviceProtocol(di)

	for protocolName, protocolStruct := range protocols {
		address := protocolStruct.(map[string]interface{})["address"].(string)
		switch protocolName {
		case Common.LORA:
			topic = "application/1/device/" + address + "/tx"
			MqttClient.Publish(topic, 0, false, encode(msg))
			break
		case Common.BLE:
			topic = "BLE"
			MqttClient.Publish(topic, 0, false, msg)
			break
		case Common.LPWA:
			topic = "lpwa/" + address + "/tx"
			MqttClient.Publish(topic, 0, false, encode(msg))
			break
		}
	}

	log.Println("==================================================[Network] Send Data : " + msg + ", Topic: " + topic)
}
