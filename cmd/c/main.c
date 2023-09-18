#define SERIAL_DEVICE_MONITOR_IMPLEMENTATION
#include "../../serial_device_monitor.h"
#include <stdio.h>

int main(int argc, char *argv[]) {
  serial_device_monitor_error err;
  serial_device_monitor_context ctx;
  err = serial_device_monitor_init(&ctx);
  if (err) {
    serial_device_monitor_deinit(&ctx);
    fprintf(stderr, "%s\n", serial_device_monitor_error_string(err));
    return 1;
  }

  struct udev_device *device;
  err = serial_device_monitor_receive(&ctx, &device);
  if (err) {
    serial_device_monitor_deinit(&ctx);
    fprintf(stderr, "%s\n", serial_device_monitor_error_string(err));
    return 1;
  }
  const char *device_path = udev_device_get_devnode(device);
  printf("device path = %s\n", device_path);

  serial_device_monitor_deinit(&ctx);
  udev_device_unref(device);
  return 0;
}
