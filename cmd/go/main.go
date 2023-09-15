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
	devpath, err := ctx.Receive()
	if err != nil {
		log.Panic(err)
	}
	log.Infof("got new terminal device at: %s", devpath)
	defer ctx.Deinit()
}
