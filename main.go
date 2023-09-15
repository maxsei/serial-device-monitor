package main

/*
#cgo pkg-config: libudev
#define SERIAL_DEVICE_MONITOR_IMPLEMENTATION
#include "src/serial_device_monitor.h"
*/
import "C"
import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	var ctx C.serial_device_monitor_context

	if err := C.serial_device_monitor_init(&ctx); err != C.serial_device_monitor_error_nil {
		C.serial_device_monitor_deinit(&ctx)
		log.Fatal(C.GoString(C.serial_device_monitor_error_string(err)))
	}

	var devicePathBuffer [4096]byte
	if err := C.serial_device_monitor_receive(&ctx, C.CBytes(devicePathBuffer[:])); err != C.serial_device_monitor_error_nil {
		C.serial_device_monitor_deinit(&ctx)
		log.Fatal(C.GoString(C.serial_device_monitor_error_string(err)))
	}
	devicePath := C.GoString(C.CBytes(devicePathBuffer[:]))
	log.Info(fmt.Sprintf("got device path: %s", devicePath))

	C.serial_device_monitor_deinit(&ctx)
}
