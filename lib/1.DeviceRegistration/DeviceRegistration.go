package DeviceRegistration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

func Request(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.InterfaceID, parameters.InterfaceID(),
		Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID())
	s, _ := json.Marshal(parameters.MicroserviceIDs())
	body := fmt.Sprintf("{%s=%v}", Parameters.MicroserviceIDs, strings.Trim(string(s), ","))
	result := header + body

	return result
}

func Response(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID(),
		Parameters.ResponseStatusCode, parameters.ResponseStatusCode())
	body := fmt.Sprintf("{%s=%s}", Parameters.DeviceID, parameters.DeviceID())
	result := header + body

	return result
}

func RQparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

	var mis []string
	for _, mi := range Common.FindFromJsonObj(data, Parameters.MicroserviceIDs).([]interface{}) {
		mis = append(mis, mi.(string))
	}

	parameters.SetMicroserviceIDs(mis)
	parameters.SetDisposableIoTRequestID(Common.FindFromJsonObj(data, Parameters.DisposableIoTRequestID).(string))
	parameters.SetInterfaceID(Common.FindFromJsonObj(data, Parameters.InterfaceID).(string))
}

func RSparsing(data []byte, parameters *Parameters.Parameter) {
}
