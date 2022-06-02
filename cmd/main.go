// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"

	disposableiot_device "github.com/keti/disposableiot-edge-gateway"
	"github.com/keti/disposableiot-edge-gateway/internal/driver"
)

const (
	serviceName string = "disposableiot-device"
)

func main() {
	d := driver.NewDisposableiotDeviceDriver()
	startup.Bootstrap(serviceName, disposableiot_device.Version, d)
}
