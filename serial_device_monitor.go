package serialdevicemonitor

/*
#cgo pkg-config: libudev
#include <libudev.h>
*/
import "C"
import "errors"

func NewContext() (*Context, error) {
	var ret Context
	ret.udev = C.udev_new()
	if ret.udev == nil {
		ret.Deinit()
		return nil, errors.New("could not create udev context")
	}
	return &ret, nil
}

type Context struct {
	udev *C.struct_udev
}

func (c *Context) Deinit() {
	if c.udev != nil {
		C.udev_unref(c.udev)
	}
}

func NewMonitor(ctx *Context) (*Monitor, error) {
	var ret Monitor
	ret.ctx = ctx
	ret.monitor = C.udev_monitor_new_from_netlink(ret.ctx.udev, C.CString("udev"))
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
	if C.udev_monitor_enable_receiving(ret.monitor) != 0 {
		ret.Deinit()
		return nil, errors.New("could not enable receiving for monitor")
	}
	return &ret, nil
}

type Monitor struct {
	ctx     *Context
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
	if m.monitor != nil {
		C.udev_monitor_unref(m.monitor)
	}
}

type Device struct {
	device *C.struct_udev_device
}

func (d *Device) DeviceNode() string {
	devnode := C.udev_device_get_devnode(d.device)
	if devnode == nil {
		return ""
	}
	return C.GoString(devnode)
}

func (d *Device) Properties() map[string]string {
	entries := C.udev_device_get_properties_list_entry(d.device)
	return udevEntriesToMap(entries)
}

func (d *Device) Deinit() {
	C.udev_device_unref(d.device)
}

func udevEntriesToMap(entries *C.struct_udev_list_entry) map[string]string {
	ret := make(map[string]string)
	for entry := entries; entry != nil; entry = C.udev_list_entry_get_next(entry) {
		k := C.udev_list_entry_get_name(entry)
		v := C.udev_list_entry_get_value(entry)
		ret[C.GoString(k)] = C.GoString(v)
	}
	return ret
}

func NewEnumerator(ctx *Context) (*Enumerator, error) {
	var ret Enumerator
	ret.ctx = ctx
	enumerate := C.udev_enumerate_new(ret.ctx.udev)
	if enumerate == nil {
		return nil, errors.New("could not enumerate devices")
	}
	ret.enumerate = enumerate
	if C.udev_enumerate_add_match_subsystem(ret.enumerate, C.CString("tty")) != 0 {
		ret.Deinit()
		return nil, errors.New("could not add match for subsystem 'tty'")
	}
	if C.udev_enumerate_scan_devices(ret.enumerate) != 0 {
		ret.Deinit()
		return nil, errors.New("enumeration failed while scannning devices")
	}
	deviceEntries := udevEntriesToMap(C.udev_enumerate_get_list_entry(enumerate))
	ret.devices = make([]string, 0, len(deviceEntries))
	for path := range deviceEntries {
		ret.devices = append(ret.devices, path)
	}
	return &ret, nil
}

type Enumerator struct {
	ctx       *Context
	enumerate *C.struct_udev_enumerate
	i         int
	devices   []string
}

func (e *Enumerator) Next() *Device {
	if e.i >= len(e.devices) {
		return nil
	}
	devpath := e.devices[e.i]
	device := C.udev_device_new_from_syspath(e.ctx.udev, C.CString(devpath))
	e.i += 1
	return &Device{device}
}

func (e *Enumerator) Deinit() {
	if e.enumerate != nil {
		C.udev_enumerate_ref(e.enumerate)
	}
}
