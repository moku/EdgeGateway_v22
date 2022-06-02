package MicroserviceRun

import (
	"bytes"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

//Request
func Request(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.InterfaceID, parameters.InterfaceID(),
		Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID())

	body := fmt.Sprintf("{%s=[%s]}", Parameters.MicroserviceIDs, strings.Join(parameters.MicroserviceIDs()[:], ","))

	result := header + body

	return result
}

//Response
func Response(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

}

func RQparsing(data []byte, parameters *Parameters.Parameter) {

	headerMap := make(map[string]string)

	//모든 trim 제거
	data = bytes.Replace(data, []byte(" "), []byte(""), -1)
	data = bytes.Replace(data, []byte("\t"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\r"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\n"), []byte(""), -1)
	fmt.Println("Data >>> " + string(data))

	temp := bytes.SplitAfterN(data, []byte("}"), 2)
	for i := range temp {
		temp[i] = bytes.ReplaceAll(temp[i], []byte("{"), []byte(""))
		temp[i] = bytes.ReplaceAll(temp[i], []byte("}"), []byte(""))

		var header [][]byte
		if i == 1 { //body

			header = bytes.Split(temp[i], []byte(";"))
			for j := range header {
				//fmt.Println(string(header[j]))
				headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			}
		} else { //header
			header = bytes.Split(temp[i], []byte(";"))
			for j := range header {
				//fmt.Println(string(header[j]))
				headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			}
		}
	}

	fmt.Println("==============================================")

	for key, value := range headerMap {
		fmt.Println("key=>"+key, "value=>"+value)
	}

	var mis []string
	tempMis := bytes.Split([]byte(headerMap["mis"]), []byte(","))
	for q := range tempMis {
		mis = append(mis, string(tempMis[q]))
	}

	parameters.SetMicroserviceIDs(mis)
	parameters.SetDisposableIoTRequestID(headerMap["dri"])

}

func RSparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

}
