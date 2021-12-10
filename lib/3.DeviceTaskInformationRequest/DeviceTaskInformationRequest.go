package DeviceTaskInformationRequest

import (
	"bytes"
	"fmt"

	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
)

//============================

type RsparsingTif struct {
	Tif []TifObject `json:"tif"`
}

type TifObject struct { /// 이걸 만들어야함
	Ti int      `json:"ti"`
	Sp SpObject `json:"sp,omitempty"` //들어올때마다 달라짐
	Fp string   `json:"fp,omitempty"`
	To bool     `json:"to,omitempty"`
}

type SpObject struct {
	value string
}

//Request
func Request(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{%s=%s;%s=%s}", Parameters.InterfaceID, parameters.InterfaceID(),
		Parameters.DisposableIoTRequestID, parameters.DisposableIoTRequestID())

	body := fmt.Sprintf("{%s=%s}", Parameters.MicroserviceID, parameters.MicroserviceID())

	result := header + body

	return result
}

//Response
func Response(parameters *Parameters.Parameter) string {

	header := fmt.Sprintf("{dri=%s;rsc=%d}", parameters.DisposableIoTRequestID(), 200)
	//fmt.Println(header)

	switch parameters.MicroserviceIDs() {

	}
	var tifObjectArray [][]string
	tifObjectArray = append(tifObjectArray, []string{"ti=1;sp={mast=200;mist=-50;};to=TRUE"})
	tifObjectArray = append(tifObjectArray, []string{"ti=2;"})

	msg := TifBulider(tifObjectArray)

	body := fmt.Sprintf("{%s}", msg)
	//fmt.Println(body)
	result := header + body
	fmt.Println(result)

	return result

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
			temp[i] = bytes.ReplaceAll(temp[i], []byte("["), []byte(""))
			temp[i] = bytes.ReplaceAll(temp[i], []byte("]"), []byte(""))
			header = bytes.Split(temp[i], []byte(";"))
			for j := range header {
				fmt.Println(string(header[j]))
				headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			}

		} else { //header
			header = bytes.Split(temp[i], []byte(";"))
			for j := range header {
				fmt.Println(string(header[j]))
				headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
			}
		}

		fmt.Println("==============================================")

		var mis []string
		tempMis := bytes.Split([]byte(headerMap["mis"]), []byte(","))
		for q := range tempMis {
			mis = append(mis, string(tempMis[q]))
		}
		parameters.SetMicroserviceIDs(mis)
		//fmt.Println(parameters.MicroserviceIDs())
		parameters.SetDisposableIoTRequestID(headerMap["dri"])

		//Response(parameters)
	}

}

func RSparsing(data map[string]interface{}, parameters *Parameters.Parameter) {

	tif := make(map[string]Parameters.TaskInformationObj)

	tifJson := Common.FindFromJsonObj(data, Parameters.TaskInformation).([]interface{})

	for _, v := range tifJson {
		element := v.(map[string]interface{})
		ti := element[Parameters.TaskID].(string)

		tif[ti] = Parameters.TaskInformationObj{TaskOrchestration: element[Parameters.TaskOrchestration].(string), StaticTaskParameter: make(map[string]string), FlexibleTaskParameter: make(map[string]string)}

		spArr := element[Parameters.StaticTaskParameter]
		if spArr != nil {
			for k, v := range spArr.(map[string]interface{}) {
				tif[ti].StaticTaskParameter[k] = v.(string)
			}
		}

		fpArr := element[Parameters.FlexibleTaskParameter]
		if fpArr != nil {
			for k, v := range fpArr.(map[string]interface{}) {
				tif[ti].FlexibleTaskParameter[k] = v.(string)
			}
		}
	}

	parameters.SetTaskInformation(tif)
}

func TifBulider(tifObjArray [][]string) string {

	fmt.Println("============Tif")
	//var objArray []string

	var result string
	for i := range tifObjArray {
		var str string
		for j := range tifObjArray[i] {
			str += fmt.Sprintf("%s", tifObjArray[i][j])
		}
		result += fmt.Sprintf("{%s},", str)
		//fmt.Println(result)
	}
	b := []byte(result)
	//fmt.Println(bytes.LastIndex(b,[]byte(",")))
	b[bytes.LastIndex(b, []byte(","))] = 0
	//fmt.Println(string(b))

	result = fmt.Sprintf("tif=[%s]", b)
	fmt.Println(result)
	return result
}
