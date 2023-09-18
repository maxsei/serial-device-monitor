package serialdevicemonitor

/*
#cgo pkg-config: libudev
#define SERIAL_DEVICE_MONITOR_IMPLEMENTATION
#include "serial_device_monitor.h"
*/
import "C"

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

func (c *Context) Receive() (*Device, error) {
	var device *C.struct_udev_device
	if err := Error(C.serial_device_monitor_receive(c.ctx, (**C.struct_udev_device)(&device))); !err.isNil() {
		return nil, err
	}
	return &Device{device}, nil
}

func (c *Context) Deinit() {
	C.serial_device_monitor_deinit(c.ctx)
}

type Device struct {
	device *C.struct_udev_device
}

func (d Device) DeviceNode() string {
	devnode := C.udev_device_get_devnode(d.device)
	if devnode == nil {
		return ""
	}
	return C.GoString(devnode)
}

func (d Device) Properties() map[string]string {
	ret := make(map[string]string)
	for entry := C.udev_device_get_properties_list_entry(d.device); entry != nil ; entry = C.udev_list_entry_get_next(entry) {
		k := C.udev_list_entry_get_name(entry);
		v := C.udev_list_entry_get_value(entry);
		ret[C.GoString(k)] = C.GoString(v)
	}
	return ret
}


func (d *Device) Deinit() {
	C.udev_device_unref(d.device)
}
