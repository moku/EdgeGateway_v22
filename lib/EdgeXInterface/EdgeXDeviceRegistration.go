package EdgeXInterface

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	EdgeXURL "github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface/URL"
)

type Address struct {
	Address string `json:"address"`
}

type ServiceInfo struct {
	Name string `json:"name"`
}

type ProfileInfo struct {
	Name string `json:"name"`
}

type DevRegParams struct {
	Name           string                 `json:"name"`
	Protocols      map[string]interface{} `json:"protocols"`
	AdminState     string                 `json:"adminState`
	OperatingState string                 `json:"operatingState`
	Service        ServiceInfo            `json:"service"`
	Profile        ProfileInfo            `json:"profile"`
	Labels         []string               `json:"labels"`
}

func NewParameter() *DevRegParams {
	res := DevRegParams{}
	/*addr := Address{"NA"}
	prtcOther := ProtocolOther{addr}
	res.Protocols = prtcOther*/
	res.AdminState = "unlocked"
	res.OperatingState = "enabled"
	return &res
}

func DeviceRegistration(param *DevRegParams) {
	tempBody, _ := json.Marshal(param)
	body := bytes.NewBuffer(tempBody)
	res, _ := http.Post(EdgeXURL.DeviceRegistration, "application/json", body)
	if res != nil {
		defer res.Body.Close()
	}
	io.Copy(ioutil.Discard, res.Body)
}

func (p *DevRegParams) SetDeviceService(name string) {
	sInfo := ServiceInfo{name}
	p.Service = sInfo
}

func (p *DevRegParams) SetDeviceProfile(name string) {
	pInfo := ProfileInfo{name}
	p.Profile = pInfo
}

/*func (p *DevRegParams) SetDeviceProtocol(protocol string, addr string) {
	if protocol == "BLE" {
		p.Protocols = Protocols{Address{addr}, Address{"NA"}}
	} else if protocol == "LORA" {
		p.Protocols = Protocols{Address{"NA"}, Address{addr}}
	}
}*/

func (p *DevRegParams) SetDeviceLabel(labels []string) {
	for _, v := range labels {
		p.Labels = append(p.Labels, v)
	}
}
