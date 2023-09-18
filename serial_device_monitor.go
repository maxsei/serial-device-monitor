package serialdevicemonitor

/*
#cgo pkg-config: libudev
#include <libudev.h>
*/
import "C"
import "errors"

func NewMonitor() (*Monitor, error) {
	var ret Monitor
	ret.udev = C.udev_new()
	if ret.udev == nil {
		ret.Deinit()
		return nil, errors.New("could not create udev context")
	}
	ret.monitor = C.udev_monitor_new_from_netlink(ret.udev, C.CString("udev"))
	if ret.monitor == nil {
		ret.Deinit()
		return nil, errors.New("could not create udev monitor")
	}

	if C.udev_monitor_filter_add_match_subsystem_devtype(
		ret.monitor, C.CString("tty"), nil,
	) != 0 {
		ret.Deinit()
		return nil, errors.New("could not initialize monitory with 'tty' subsystem filter")
	}

  if (C.udev_monitor_enable_receiving(ret.monitor) != 0) {
		ret.Deinit()
    return nil, errors.New("could not enable receiving for monitor")
  }
	return &ret, nil
}

type Monitor struct {
	udev    *C.struct_udev
	monitor *C.struct_udev_monitor
}

func (m *Monitor) Receive() (*Device, error) {
	for true {
		device := C.udev_monitor_receive_device(m.monitor)
		if device == nil {
			continue
		}
		action := C.GoString(C.udev_device_get_action(device))
		if action != "add" {
			C.udev_device_unref(device)
			continue
		}
		return &Device{device}, nil
	}
	panic("unreachable")
}

func (m *Monitor) Deinit() {
	if m.udev != nil {
		C.udev_unref(m.udev)
	}
	if m.monitor != nil {
		C.udev_monitor_unref(m.monitor)
	}
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
	for entry := C.udev_device_get_properties_list_entry(d.device); entry != nil; entry = C.udev_list_entry_get_next(entry) {
		k := C.udev_list_entry_get_name(entry)
		v := C.udev_list_entry_get_value(entry)
		ret[C.GoString(k)] = C.GoString(v)
	}
	return ret
}

func (d *Device) Deinit() {
	C.udev_device_unref(d.device)
}
