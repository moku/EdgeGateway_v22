package EdgeXInterface

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface/URL"
)

type DevRegInfo struct {
	DeviceID string `json:"device_id"`
}

type Notification struct {
	Sender   string   `json:"sender"`
	Category string   `json:category`
	Severity string   `json: severity`
	Content  string   `json:"content"`
	Labels   []string `json:"labels"`
}

type NotificationParams struct {
	ApiVersion   string       `json:apiVersion`
	Notification Notification `json:notification`
}

func DeleteNotification(slug string) {
	req, _ := http.NewRequest("DELETE", URL.DeleteNotification+"/"+slug, nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	io.Copy(ioutil.Discard, resp.Body)
}

func SendNotification(deviceName string) {

	DeleteNotification(deviceName)

	devRegInfo := DevRegInfo{deviceName}
	content, _ := json.Marshal(devRegInfo)
	notificationParams := make([]NotificationParams, 1)
	notificationParams[0] = NotificationParams{"v2", Notification{"disposableiot-gateway", "DisposableIoT", "NORMAL", string(content), []string{"device-registration"}}}

	str, _ := json.Marshal(notificationParams)
	body := bytes.NewBufferString(string(str))

	resp, err := http.Post(URL.InsertNotification, "application/json", body)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	io.Copy(ioutil.Discard, resp.Body)
}
