package MicroserviceInputParameterSet

import (
	"bytes"
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

//Request
func Request(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.InterfaceID, parameters.InterfaceID(),
		Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID())

	ip := ""
	for k, v := range parameters.InputParameter() {
		ip += fmt.Sprintf("%s=%s;", k, v)
	}
	ip = ip[:len(ip)-1] // remove the last semicolon
	body := fmt.Sprintf("{%s=%s;%s={%s}}", Parameters.MicroserviceID, parameters.MicroserviceID(),
		Parameters.MicroserviceInputParameter, ip)

	result := header + body

	return result
}

//Response
func Response(parameters *Parameters.Parameter, mqttClient mqtt.Client) {

}

func RQparsing(data []byte, parameters *Parameters.Parameter) {

	headerMap := make(map[string]string)
	//var sfs []reflect.StructField
	//objArray := make(map[int]string)

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

			//header = bytes.Split(temp[i], []byte(";"))
			//for j := range header {
			//	//fmt.Println(string(header[j]))
			//	headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			//}

		} else { //header
			header = bytes.Split(temp[i], []byte(";"))
			for j := range header {
				//fmt.Println(string(header[j]))
				headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			}
		}

		fmt.Println("==============================================")

		for key, value := range headerMap {
			fmt.Println("key=>"+key, "value=>"+value)
		}

		var mis []int
		tempMis := bytes.Split([]byte(headerMap["mis"]), []byte(","))
		for q := range tempMis {
			byteToInt, _ := strconv.Atoi(string(tempMis[q]))
			mis = append(mis, byteToInt)
		}

		//Response Message
		//n,e := strconv.Atoi(headerMap["dri"])
		//fmt.Println(e)
		//n++
		//parameters.SetDisposableIoTRequestID(strconv.Itoa(n))
		parameters.SetDisposableIoTRequestID(headerMap["dri"])

		//	Response(parameters)

		//fmt.Println(string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[0]),[]byte(";"),[]byte("")),[]byte("="))[1]))

	}

}

func RSparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

}
