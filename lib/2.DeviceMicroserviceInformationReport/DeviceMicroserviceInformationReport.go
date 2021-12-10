package DeviceMicroserviceInformationReport

import (
	"fmt"

	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

func Request(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{if=%s;dri=%s;di=%s}", parameters.InterfaceID(), parameters.DisposableIoTRequestID(), parameters.DeviceID())
	body := fmt.Sprintf("{%s}", parameters.MicroserviceInformation())
	result := header + body

	return result
}

func Response(parameters *Parameters.Parameter) string {

	result := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID(),
		Parameters.ResponseStatusCode, parameters.ResponseStatusCode())

	return result
}

func RQparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

	var mif map[string][]string = make(map[string][]string)

	mifJson := Common.FindFromJsonObj(data, Parameters.MicroserviceInformation).([]interface{})
	for _, v := range mifJson {
		element := v.(map[string]interface{})
		mi := element[Parameters.MicroserviceID].(string)
		opsArr := element[Parameters.MicroserviceOutputParameters].([]interface{})

		for _, ops := range opsArr {
			mif[mi] = append(mif[mi], ops.(string))
		}
	}

	parameters.SetDisposableIoTRequestID(Common.FindFromJsonObj(data, Parameters.DisposableIoTRequestID).(string))
	parameters.SetInterfaceID(Common.FindFromJsonObj(data, Parameters.InterfaceID).(string))
	parameters.SetDeviceID(Common.FindFromJsonObj(data, Parameters.DeviceID).(string))
	parameters.SetMicroserviceInformation(mif)
}

func RSparsing(data []byte, parameters *Parameters.Parameter) {

}
