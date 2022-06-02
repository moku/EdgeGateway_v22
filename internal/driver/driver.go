// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a implementation of a ProtocolDriver interface.
//
package driver

import (
	list2 "container/list"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"

	MicroserviceCreation "github.com/keti/disposableiot-edge-gateway/lib/10.MicroserviceCreation"
	MicroserviceRun "github.com/keti/disposableiot-edge-gateway/lib/11.MicroserviceRun"
	MicroserviceInputParameterSet "github.com/keti/disposableiot-edge-gateway/lib/12.MicroserviceInputParameterSet"
	MicroserviceOutputParameterRead "github.com/keti/disposableiot-edge-gateway/lib/13.MicroserviceOutputParameterRead"
	MicroserviceStop "github.com/keti/disposableiot-edge-gateway/lib/15.MicroserviceStop"
	MicroserviceDelete "github.com/keti/disposableiot-edge-gateway/lib/16.MicroserviceDelete"
	TaskRun "github.com/keti/disposableiot-edge-gateway/lib/20.TaskRun"
	TaskParameterSet "github.com/keti/disposableiot-edge-gateway/lib/21.TaskParameterSet"
	TaskParameterRead "github.com/keti/disposableiot-edge-gateway/lib/22.TaskParameterRead"
	TaskParameterReadAll "github.com/keti/disposableiot-edge-gateway/lib/23.TaskParameterReadAll"
	TaskStop "github.com/keti/disposableiot-edge-gateway/lib/24.TaskStop"
	DeviceTaskInformationRequest "github.com/keti/disposableiot-edge-gateway/lib/3.DeviceTaskInformationRequest"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/keti/disposableiot-edge-gateway/lib/Parameters"
	"github.com/keti/disposableiot-edge-gateway/lib/ResourceName"

	"github.com/keti/disposableiot-edge-gateway/lib/Common"

	_ "modernc.org/ql/driver"
)

var once sync.Once
var driver *DisposableiotDeviceDriver
var sdkService sdk.DeviceService
var logging = true
var MqttClient mqtt.Client
var globalDri int
var DriList []string
var LastNode string
var parametersHistory = make(map[string]*Parameters.Parameter) // [Dri]Parameter
var TaskList []*list2.Element
var ListPointer *list2.Element

var MqttAddr = "localhost:1883"

type DisposableiotDeviceDriver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
}

const (
	Comma      = "%2C"
	LeftBrace  = "%7B"
	RightBrace = "%7D"
	Equal      = "%3D"
)

func UpdateDri(dri int) {
	if globalDri < dri {
		globalDri = dri
	}
}

func NewDisposableiotDeviceDriver() dsModels.ProtocolDriver {
	once.Do(func() {
		driver = new(DisposableiotDeviceDriver)
	})
	return driver
}

func (d *DisposableiotDeviceDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	d.lc = lc
	d.asyncCh = asyncCh
	sdkService = *sdk.RunningService()

	if !logging {
		log.SetOutput(ioutil.Discard)
	}

	log.Println("==================================================Initialize")

	opts := mqtt.NewClientOptions().AddBroker(MqttAddr) // MQTT Server

	var MsgHandler mqtt.MessageHandler = MsgHandler
	opts.SetDefaultPublishHandler(MsgHandler)
	MqttClient = mqtt.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	var topicList = make(map[string]byte)

	topicList["pipeLineDown"] = 0
	topicList["registrationUp"] = 0
	topicList["messageUp"] = 0
	topicList["BLE"] = 0
	topicList["application/#"] = 0
	topicList[Common.LPWA+"/#"] = 0

	if token := MqttClient.SubscribeMultiple(topicList, nil); token.Wait() && token.Error() != nil {
		//fmt.Println()
		log.Println(token.Error())
		return nil
	}

	return nil
}

func (d *DisposableiotDeviceDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	for _, req := range reqs {
		for k := range req.Attributes {
			str := req.Attributes[k].(string)
			str = strings.ReplaceAll(str, LeftBrace, "{")
			str = strings.ReplaceAll(str, RightBrace, "}")
			str = strings.ReplaceAll(str, Equal, "=")
			str = strings.ReplaceAll(str, Comma, ",")
			req.Attributes[k] = str

			fmt.Println("[READ]" + req.DeviceResourceName + ", " + k + " : " + str)
		}
	}

	res = make([]*dsModels.CommandValue, len(reqs))

	for i, req := range reqs {

		globalDri++
		strDri := strconv.Itoa(globalDri)
		parameters := parametersHistory[strDri]
		parameters = ParseCmd(reqs[i].Attributes["urlRawQuery"].(string))
		parameters.SetDeviceID(deviceName)
		parameters.SetDisposableIoTRequestID(strDri)

		switch req.DeviceResourceName {
		case "dti": //DeviceTaskInformation  //3

			log.Println("==================================================[CMD] DeviceTaskInformation")

			parameters.SetInterfaceID(ResourceName.DeviceTaskInformationRequest)
			msg := DeviceTaskInformationRequest.Request(parameters)
			SendDownlinkMsg(msg, deviceName)

			res[i], _ = dsModels.NewCommandValue(req.DeviceResourceName, "string", "dti")

			break
		case "mopr": //MicroserviceOutputParameterRead 13

			log.Println("==================================================[CMD] MicroserviceOutputParameterRead")

			parameters.SetInterfaceID(ResourceName.MicroserviceOutputParameterRead)
			msg := MicroserviceOutputParameterRead.Request(parameters)
			SendDownlinkMsg(msg, deviceName)

			res[i], _ = dsModels.NewCommandValue(req.DeviceResourceName, "string", "mopr")

			break
		case "tpr": //TaskParameterRead 22
			log.Println("==================================================[CMD] TaskParameterRead")

			parameters.SetInterfaceID(ResourceName.TaskParameterRead)
			msg := TaskParameterRead.Request(parameters)
			SendDownlinkMsg(msg, deviceName)

			res[i], _ = dsModels.NewCommandValue(req.DeviceResourceName, "string", "tpr")

			break
		case "tpra": //TaskParameterReadAll 23
			log.Println("==================================================[CMD] TaskParameterReadAll")

			parameters.SetInterfaceID(ResourceName.TaskParameterReadAll)
			msg := TaskParameterReadAll.Request(parameters)
			SendDownlinkMsg(msg, deviceName)

			res[i], _ = dsModels.NewCommandValue(req.DeviceResourceName, "string", "tpra")

			break
		}

		parametersHistory[strDri] = parameters
	}

	return res, nil
}

func (d *DisposableiotDeviceDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	log.Println(reqs[0].Attributes)

	for _, cmd := range params {
		data, _ := cmd.StringValue()
		fmt.Println("[WRITE]" + cmd.DeviceResourceName + ", " + data)
	}

	for k, v := range params {

		globalDri++
		strDri := strconv.Itoa(globalDri)
		//var parameters *Parameters.Parameter = parametersHistory[strDri]
		parameters := Parameters.NewParameter()
		str, _ := params[k].StringValue()
		parameters = ParseCmd(str)
		parameters.SetDeviceID(deviceName)
		parameters.SetDisposableIoTRequestID(strDri)

		switch v.DeviceResourceName {
		case "mc": //MicroserviceCreation //10
			log.Println("==================================================[CMD] MicroserviceCreation")

			parameters.SetInterfaceID(ResourceName.MicroserviceCreation)
			msg := MicroserviceCreation.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "mr": //MicroserviceRun 11
			log.Println("==================================================[CMD] MicroserviceRun")

			parameters.SetInterfaceID(ResourceName.MicroserviceRun)
			msg := MicroserviceRun.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "mips": //MicroserviceInputParameterSet 12
			log.Println("==================================================[CMD] MicroserviceInputParameterSet")

			parameters.SetInterfaceID(ResourceName.MicroserviceInputParameterSet)
			msg := MicroserviceInputParameterSet.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "ms": //MicroserviceStop  15
			log.Println("==================================================[CMD] MicroserviceStop")

			parameters.SetInterfaceID(ResourceName.MicroserviceStop)
			msg := MicroserviceStop.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "md": //MicroserviceDelete 16
			log.Println("==================================================[CMD] MicroserviceDelete")

			parameters.SetInterfaceID(ResourceName.MicroserviceDelete)
			msg := MicroserviceDelete.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "ton": //TaskOn  20
			log.Println("==================================================[CMD] TaskOn")

			parameters.SetInterfaceID(ResourceName.TaskOn)
			msg := TaskRun.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break

		case "tps": //TaskParameterSet 21
			log.Println("==================================================[CMD] TaskParameterSet")

			parameters.SetInterfaceID(ResourceName.TaskParameterSet)
			msg := TaskParameterSet.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break
		case "tof": //TaskOff
			log.Println("==================================================[CMD] TaskOff")

			parameters.SetInterfaceID(ResourceName.TaskOff)
			msg := TaskStop.Request(parameters)

			SendDownlinkMsg(msg, deviceName)

			break
		}

		parametersHistory[strDri] = parameters
	}

	return nil
}

func (d *DisposableiotDeviceDriver) Stop(force bool) error {
	d.lc.Info("DisposableDeviceDriver.Stop: device-virtual driver is stopping...")
	return nil
}

func (d *DisposableiotDeviceDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("a new Device is added: %s", deviceName)
	return nil
}

func (d *DisposableiotDeviceDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("Device %s is updated", deviceName)
	return nil
}

func (d *DisposableiotDeviceDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Debugf("Device %s is removed", deviceName)
	return nil
}
