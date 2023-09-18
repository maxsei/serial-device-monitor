package main

import (
	serialdevicemonitor "github.com/maxsei/serial-device-monitor"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	ctx, err := serialdevicemonitor.Init()
	if err != nil {
		log.Panic(err)
	}
	device, err := ctx.Receive()
	if err != nil {
		log.Panic(err)
	}
	defer device.Deinit()
	devpath := device.DeviceNode()
	log.Infof("got new terminal device at: %s", devpath)
	props := device.Properties()
	for k, v := range props {
		log.Infof("%s: %s", k, v)
	}
	defer ctx.Deinit()
}
