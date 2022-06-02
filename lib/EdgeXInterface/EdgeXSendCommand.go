package EdgeXInterface

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface/URL"
)

func SendCommand(device string, command string, msg string) {

	client := &http.Client{}
	body := bytes.NewBufferString(string(msg))

	req, err := http.NewRequest(http.MethodPut, URL.SendDeviceCommand+"/"+device+"/"+command, body)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	io.Copy(ioutil.Discard, resp.Body)
}
