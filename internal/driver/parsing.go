package driver

import (
	"bytes"
	list2 "container/list"
	"encoding/json"
	"regexp"
	"time"

	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	DeviceRegistration "github.com/keti/disposableiot-edge-gateway/lib/1.DeviceRegistration"
	MicroserviceCreation "github.com/keti/disposableiot-edge-gateway/lib/10.MicroserviceCreation"
	MicroserviceRun "github.com/keti/disposableiot-edge-gateway/lib/11.MicroserviceRun"
	MicroserviceInputParameterSet "github.com/keti/disposableiot-edge-gateway/lib/12.MicroserviceInputParameterSet"
	MicroserviceOutputParameterRead "github.com/keti/disposableiot-edge-gateway/lib/13.MicroserviceOutputParameterRead"
	MicroserviceOutputReport "github.com/keti/disposableiot-edge-gateway/lib/14.MicroserviceOutputReport"
	MicroserviceStop "github.com/keti/disposableiot-edge-gateway/lib/15.MicroserviceStop"
	MicroserviceDelete "github.com/keti/disposableiot-edge-gateway/lib/16.MicroserviceDelete"
	DeviceMicroserviceInformationReport "github.com/keti/disposableiot-edge-gateway/lib/2.DeviceMicroserviceInformationReport"
	TaskRun "github.com/keti/disposableiot-edge-gateway/lib/20.TaskRun"
	TaskParameterSet "github.com/keti/disposableiot-edge-gateway/lib/21.TaskParameterSet"
	TaskParameterRead "github.com/keti/disposableiot-edge-gateway/lib/22.TaskParameterRead"
	TaskParameterReadAll "github.com/keti/disposableiot-edge-gateway/lib/23.TaskParameterReadAll"
	TaskOff "github.com/keti/disposableiot-edge-gateway/lib/24.TaskStop"
	DeviceTaskInformationRequest "github.com/keti/disposableiot-edge-gateway/lib/3.DeviceTaskInformationRequest"
	"github.com/keti/disposableiot-edge-gateway/lib/Common"
	"github.com/keti/disposableiot-edge-gateway/lib/EdgeXInterface"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
	"github.com/keti/disposableiot-edge-gateway/lib/ResourceName"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func WrappingWithQuote(str string) string {
	return "\"" + str + "\""
}

func MsgToJson(msg string) map[string]interface{} {

	msgJsonObj := make(map[string]interface{})

	if strings.Contains(msg, "}{") == true {
		// msg contains both header and body

		msgTemp := strings.Replace(msg, "}", "},body:", 1) // split body section
		msgTemp = strings.ReplaceAll(msgTemp, ";", ",")
		msgTemp = strings.ReplaceAll(msgTemp, "=", ":")

		r := regexp.MustCompile("[^\\[{:,}\\]]+")                        // find key and value
		msgJsonStr := r.ReplaceAllStringFunc(msgTemp, WrappingWithQuote) // wrap key and value with ""
		msgJsonStr = "{\"header\":" + msgJsonStr + "}"

		log.Println(msgJsonStr)

		err := json.Unmarshal([]byte(msgJsonStr), &msgJsonObj)

		if err != nil {
			panic(err)
		}
	} else {
		// msg contains only header
		msgTemp := strings.ReplaceAll(msg, ";", ",")
		msgTemp = strings.ReplaceAll(msgTemp, "=", ":")

		r := regexp.MustCompile("[^\\[{:,}\\]]+")                        // find key and value
		msgJsonStr := r.ReplaceAllStringFunc(msgTemp, WrappingWithQuote) // wrap key and value with ""

		err := json.Unmarshal([]byte(msgJsonStr), &msgJsonObj)

		if err != nil {
			panic(err)
		}
	}

	return msgJsonObj
}

func ParseCmd(rawParams string) *Parameters.Parameter {

	var tempMap = map[string]string{}
	temps := strings.Split(rawParams, "&")
	for i := range temps {
		temp := strings.SplitN(temps[i], "=", 2)
		tempMap[temp[0]] = temp[1]
	}

	parameters := Parameters.NewParameter()

	for k, v := range tempMap {
		switch k {
		case "mi":
			parameters.SetMicroserviceID(v)

			break
		case "mis":
			tempSplit := strings.Split(v, ",")
			var mis []string
			for _, mi := range tempSplit {
				mis = append(mis, mi)
			}
			parameters.SetMicroserviceIDs(mis)

			break
		case "ti":
			parameters.SetTaskID(v)

			break
		case "tis":
			tempSplit := strings.Split(v, ",")
			var tis []string
			for _, ti := range tempSplit {
				tis = append(tis, ti)
			}
			parameters.SetTaskIDs(tis)

			break
		case "ip":
			// ip = {k1=v1;k2=v2;...}
			tempStr := strings.ReplaceAll(v, "{", "")
			tempStr = strings.ReplaceAll(tempStr, "}", "")
			ipObj := strings.Split(tempStr, ",")

			ip := make(map[string]string)
			for i := range ipObj {
				ipPair := strings.Split(ipObj[i], "=")
				ip[ipPair[0]] = ipPair[1]
			}
			parameters.SetInputParameter(ip)

			break
		case "fp":
			// fp = {k1=v1;k2=v2;...}
			tempStr := strings.ReplaceAll(v, "{", "")
			tempStr = strings.ReplaceAll(tempStr, "}", "")
			fpObj := strings.Split(tempStr, ",")

			fp := make(map[string]int)
			for _, fpObj := range fpObj {
				fpPair := strings.Split(fpObj, "=")
				fp[fpPair[0]], _ = strconv.Atoi(fpPair[1])
			}
			parameters.SetFlexibleTaskParameter(fp)

			break
		case "fps":
			tempSplit := strings.Split(v, ",")
			var fps []string
			fps = append(fps, tempSplit[0])
			for _, fp := range tempSplit[1:] {
				fps = append(fps, fp)
			}
			parameters.SetFlexibleTaskParameters(fps)

			break
		case "sps":
			tempSplit := strings.Split(v, ",")
			var sps []string
			sps = append(sps, tempSplit[0])
			for _, sp := range tempSplit[1:] {
				sps = append(sps, sp)
			}
			parameters.SetStaticTaskParameters(sps)

			break
		}
	}

	return parameters
}

func ParseRequestMsg(data []byte, client mqtt.Client, protocol string, addr string) {

	log.Println("[Parsing-REQ] Get Request MSG from Device")

	msgMap := MsgToJson(string(data))

	log.Printf("[Parsing-REQ] %s\n", msgMap)

	parameters := Parameters.NewParameter()
	parameters.SetProtocol(protocol)
	parameters.SetProtocolAddr(addr)

	switch Common.FindFromJsonObj(msgMap, Parameters.InterfaceID).(string) {

	case ResourceName.DeviceRegistration: //1
		log.Println("[Parsing-REQ] DeviceRegistration")

		DeviceRegistration.RQparsing(msgMap, parameters)

		deviceID := ""
		if res, previousDeviceID := EdgeXInterface.IsRedundantDevice(parameters.ProtocolAddr()); res == true {
			// redundant registration request
			deviceID = previousDeviceID
		} else {
			deviceID, _ = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)

			device := models.Device{}
			device.AdminState = "UNLOCKED"
			device.OperatingState = "UP"
			device.Name = deviceID
			device.ServiceName = "disposable-device"
			device.ProfileName = "DisposableIoT-Device"
			device.Labels = parameters.MicroserviceIDs()
			device.Protocols = make(map[string]models.ProtocolProperties)
			device.Protocols[parameters.Protocol()] = models.ProtocolProperties{"address": parameters.ProtocolAddr()}

			_, err := sdkService.AddDevice(device)
			if err != nil {
				panic(err)
			}

			/*devRegParam := EdgeXInterface.NewParameter()
			devRegParam.Name = deviceID
			devRegParam.SetDeviceProfile("Disposable Device")
			devRegParam.SetDeviceService("disposable-device")
			devRegParam.SetDeviceProtocol(parameters.Protocol(), parameters.ProtocolAddr())
			devRegParam.SetDeviceLabel(parameters.MicroserviceIDs())
			EdgeXInterface.DeviceRegistration(devRegParam)*/
		}

		EdgeXInterface.SendNotification(deviceID)

		parameters.SetDeviceID(deviceID)

		dri, _ := strconv.Atoi(parameters.DisposableIoTRequestID())
		UpdateDri(dri)

		parameters.SetResponseStatusCode("200")
		res := DeviceRegistration.Response(parameters)

		SendDownlinkMsg(res, parameters.DeviceID())

		break
	case ResourceName.DeviceMicroserviceInformationReport: //2
		log.Println("[Parsing-REQ] DeviceMicroserviceInformationReport")

		DeviceMicroserviceInformationReport.RQparsing(msgMap, parameters)

		dri, _ := strconv.Atoi(parameters.DisposableIoTRequestID())
		UpdateDri(dri)

		parameters.SetResponseStatusCode("200")
		res := DeviceMicroserviceInformationReport.Response(parameters)

		SendDownlinkMsg(res, parameters.DeviceID())

		time.Sleep(time.Second * 10)

		msg_json := map[string]string{"mc": "mis=101"}
		msg_str, _ := json.Marshal(msg_json)
		EdgeXInterface.SendCommand(parameters.DeviceID(), "MicroserviceCreation", string(msg_str))

		/*time.Sleep(1)

		header := fmt.Sprintf("{%s=%s;%s=%d}", Parameters.InterfaceID, ResourceName.MicroserviceCreation,
			Parameters.DisposableIoTRequestID, globalDri)
		UpdateDri(globalDri)

		body := fmt.Sprintf("{%s=[%s]}", Parameters.MicroserviceIDs, "101")

		result := header + body

		SendDownlinkMsg(result, parameters.DeviceID())*/

		break
	case ResourceName.MicroserviceOutputReport: //Device ID

		log.Println("[Parsing-REQ] MicroserviceOutputReport")

		MicroserviceOutputReport.RQparsing(msgMap, parameters)

		for op, value := range parameters.OutputParameter() {
			EdgeXInterface.InsertSensingData(parameters.DeviceID(), parameters.MicroserviceID(), op, value)
		}

		dri, _ := strconv.Atoi(parameters.DisposableIoTRequestID())
		UpdateDri(dri)

		parameters.SetResponseStatusCode("200")
		res := MicroserviceOutputReport.Response(parameters)

		SendDownlinkMsg(res, parameters.DeviceID())

		if len(TaskList) > 0 {
			//Scheduler(parameters.GetData(), nil, nil)
		}

		break
	}

	//parameters.SetDriIf(driIf)
}

func ParseResponseMsg(data []byte) {

	msgMap := MsgToJson(string(data))

	log.Printf("[Parsing-RES] %s\n", msgMap)

	dri := Common.FindFromJsonObj(msgMap, Parameters.DisposableIoTRequestID).(string)
	status := Common.FindFromJsonObj(msgMap, Parameters.ResponseStatusCode).(string)

	if status != "200" {
		log.Println("[Parsing-RES] Response Code is " + status + "!!")
		delete(parametersHistory, dri)
		return
	}

	parameters := parametersHistory[dri]
	ifnum := parameters.InterfaceID()

	switch ifnum {

	case ResourceName.DeviceTaskInformationRequest: //3
		log.Println("[Parsing-RES] DeviceTaskInformationRequest")

		DeviceTaskInformationRequest.RSparsing(msgMap, parameters)
		log.Println(parameters.TaskInformation())

		break

	case ResourceName.MicroserviceCreation: //MicroserviceCreation //10
		log.Println("[Parsing-RES] MicroserviceCreation")
		MicroserviceCreation.RSparsing(msgMap, parameters)

		time.Sleep(time.Second * 10)

		msg_json := map[string]string{"mr": "mis=101"}
		msg_str, _ := json.Marshal(msg_json)
		EdgeXInterface.SendCommand(parameters.DeviceID(), "MicroserviceRun", string(msg_str))

		/*header := fmt.Sprintf("{%s=%s;%s=%d}", Parameters.InterfaceID, ResourceName.MicroserviceRun,
			Parameters.DisposableIoTRequestID, globalDri)
		UpdateDri(globalDri)

		body := fmt.Sprintf("{%s=[%s]}", Parameters.MicroserviceIDs, "101")

		result := header + body

		SendDownlinkMsg(result, parameters.DeviceID())*/

		break

	case ResourceName.MicroserviceRun: //MicroserviceRun 11
		log.Println("[Parsing-RES] MicroserviceRun")
		MicroserviceRun.RSparsing(msgMap, parameters)

		break

	case ResourceName.MicroserviceInputParameterSet: //MicroserviceInputParameterSet 12
		log.Println("[Parsing-RES] MicroserviceInputParameterSet")
		MicroserviceInputParameterSet.RSparsing(msgMap, parameters)

		break

	case ResourceName.MicroserviceOutputParameterRead: //MicroserviceOutputParameterRead 13
		log.Println("[Parsing-RES] MicroserviceOutputParameterRead")
		MicroserviceOutputParameterRead.RSparsing(msgMap, parameters)

		for op, value := range parameters.OutputParameter() {
			EdgeXInterface.InsertSensingData(parameters.DeviceID(), parameters.MicroserviceID(), op, value)
		}

		break

	case ResourceName.MicroserviceStop: //MicroserviceStop  15
		log.Println("[Parsing-RES] MicroserviceStop")
		MicroserviceStop.RSparsing(msgMap, parameters)

		break

	case ResourceName.MicroserviceDelete: //MicroserviceDelete 16
		log.Println("[Parsing-RES] MicroserviceDelete")
		MicroserviceDelete.RSparsing(msgMap, parameters)

		break

	case ResourceName.TaskOn: //TaskOn  20
		log.Println("[Parsing-RES] TaskOn")
		TaskRun.RSparsing(msgMap, parameters)

		break

	case ResourceName.TaskParameterSet: //TaskParameterSet 21
		log.Println("[Parsing-RES] TaskParameterSet")
		TaskParameterSet.RSparsing(msgMap, parameters)

		break

	case ResourceName.TaskParameterRead: //TaskParameterRead 22
		log.Println("[Parsing-RES] TaskParameterRead")
		TaskParameterRead.RSparsing(msgMap, parameters)

		break

	case ResourceName.TaskParameterReadAll: //TaskParameterReadAll 23
		log.Println("[Parsing-RES] TaskParameterReadAll")
		TaskParameterReadAll.RSparsing(msgMap, parameters)

		break

	case ResourceName.TaskOff: //TaskOff  24
		log.Println("[Parsing-RES] TaskOff")
		TaskOff.RSparsing(msgMap, parameters)

		break
	}

	if ifnum != ResourceName.MicroserviceOutputParameterRead {
		delete(parametersHistory, dri)
	}
}

func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveTask(slice []*list2.Element, s int) []*list2.Element {
	return append(slice[:s], slice[s+1:]...)
}

func resParsingData(data []byte) map[string]string {

	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(err)
		}
	}()

	headerMap := make(map[string]string)

	//모든 trim 제거
	data = bytes.Replace(data, []byte(" "), []byte(""), -1)
	data = bytes.Replace(data, []byte("\t"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\r"), []byte(""), -1)
	data = bytes.Replace(data, []byte("\n"), []byte(""), -1)

	data = bytes.ReplaceAll(data, []byte("{"), []byte(""))
	data = bytes.ReplaceAll(data, []byte("}"), []byte(""))
	//
	var header [][]byte
	header = bytes.Split(data, []byte(";"))
	for j := range header {

		headerMap[string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[0])] = string(bytes.Split(bytes.ReplaceAll(bytes.TrimSpace(header[j]), []byte(";"), []byte("")), []byte("="))[1])
	}
	//}
	//fmt.Println("rsif: "+ headerMap["if"])
	//fmt.Println("rsdri: "+ headerMap["dri"] )

	return headerMap
}
