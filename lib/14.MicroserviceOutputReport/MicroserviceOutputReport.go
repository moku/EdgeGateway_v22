package MicroserviceOutputReport

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

//Request
func Request(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

}

//Response
func Response(parameters *Parameters.Parameter) string {

	result := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID(),
		Parameters.ResponseStatusCode, parameters.ResponseStatusCode())

	return result
}

func RQparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

	parameters.SetDeviceID(Common.FindFromJsonObj(data, Parameters.DeviceID).(string))
	parameters.SetDisposableIoTRequestID(Common.FindFromJsonObj(data, Parameters.DisposableIoTRequestID).(string))
	parameters.SetMicroserviceID(Common.FindFromJsonObj(data, Parameters.MicroserviceID).(string))
	opObj := make(map[string]string)
	opJson := Common.FindFromJsonObj(data, Parameters.MicroserviceOutputParameter).(map[string]interface{})
	for k, v := range opJson {
		opObj[k], _ = v.(string)
	}
	parameters.SetOutputParameter(opObj)
}

func RSparsing(data []byte, parameters *Parameters.Parameter) {

}
