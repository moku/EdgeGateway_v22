package EdgeXInterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface/URL"
)

var MqttClient mqtt.Client = nil

type Reading struct {
	DeviceName   string `json:"deviceName"`
	ProfileName  string `json:"profileName"`
	ResourceName string `json:"resourceName"`
	Id           string `json:"id"`
	Origin       int64  `json:"origin"`
	ValueType    string `json:"valueType"`
	Value        string `json:"value"`
}

type Event struct {
	ApiVersion  string            `json:"apiVersion"`
	DeviceName  string            `json:"deviceName"`
	ProfileName string            `json:"profileName"`
	SourceName  string            `json:"sourceName"`
	Id          string            `json:"id"`
	Origin      int64             `json:"origin"`
	Readings    []Reading         `json:"readings"`
	Tags        map[string]string `json:"tags"`
}

type EventStructure struct {
	ApiVersion string `json:"apiVersion"`
	Event      Event  `json:"event"`
}

func DeliveryData(data string) {
	log.Println("==================================================[Network] Delivery Data: " + data)

	if MqttClient == nil {
		opts := mqtt.NewClientOptions().AddBroker("localhost:1883")
		MqttClient = mqtt.NewClient(opts)
		if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	MqttClient.Publish("EdgeXEvents", 0, false, data)
}

func InsertSensingData(deviceName string, microservice string, dataType string, data string) {

	id, _ := uuid.NewUUID()
	origin := time.Now().UnixNano()

	id, _ = uuid.NewUUID()
	reading := Reading{deviceName, "DisposableIoT-Device", dataType, id.String(), origin, "String", data}
	event := Event{"v2", deviceName, "DisposableIoT-Device", dataType, id.String(), origin, []Reading{reading}, map[string]string{"mi": microservice}}
	eventStructure := EventStructure{"v2", event}

	str, _ := json.Marshal(eventStructure)
	body := bytes.NewBufferString(string(str))

	fmt.Println(body)

	resp, err := http.Post(URL.InsertSensingData+"/DisposableIoT-Device/"+deviceName+"/"+dataType, "application/json", body)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	io.Copy(ioutil.Discard, resp.Body)

	//DeliveryData(string(str))
}
