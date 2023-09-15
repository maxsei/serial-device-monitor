package serialdevicemonitor

/*
#cgo pkg-config: libudev
#define SERIAL_DEVICE_MONITOR_IMPLEMENTATION
#include "serial_device_monitor.h"
*/
import "C"
import (
	"bytes"
	"io"
)

type Error C.serial_device_monitor_error

func (e Error) Error() string {
	cerr := C.serial_device_monitor_error(e)
	return C.GoString(C.serial_device_monitor_error_string(cerr))
}

func (e Error) isNil() bool {
	return int(e) == C.serial_device_monitor_error_nil
}

func Init() (*Context, error) {
	var ctx C.serial_device_monitor_context
	ret := &Context{ctx: &ctx}
	if err := Error(C.serial_device_monitor_init(&ctx)); !err.isNil() {
		ret.Deinit()
		return nil, err
	}
	return ret, nil
}

type Context struct {
	ctx *C.serial_device_monitor_context
}

func (c *Context) Receive() (string, error) {
	var devicePathBuffer [4096]byte
	if err := Error(C.serial_device_monitor_receive(
		c.ctx,
		(*C.char)(C.CBytes(devicePathBuffer[:]))),
	); !err.isNil() {
		return "", err
	}
	end := bytes.IndexByte(devicePathBuffer[:], 0)
	ret := bytes.NewBuffer(nil)
	if _, err := io.Copy(ret, bytes.NewBuffer(devicePathBuffer[:end])); err != nil {
		return "", err
	}
	return ret.String(), nil
}

func (c *Context) Deinit() {
	C.serial_device_monitor_deinit(c.ctx)
}
