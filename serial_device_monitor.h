#ifndef INCLUDE_SERIAL_DEVICE_MONITOR_H
#define INCLUDE_SERIAL_DEVICE_MONITOR_H

#include <fcntl.h>
#include <libudev.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef enum {
  serial_device_monitor_error_nil,
  serial_device_monitor_error_null_device,
  serial_device_monitor_error_udev_ctx_init,
  serial_device_monitor_error_udev_mon_init,
  serial_device_monitor_error_add_subsys_filter,
  serial_device_monitor_error_mon_enable_failed,
  /* serial_device_monitor_error_device_recieve_failed */
} serial_device_monitor_error;
extern const char *
serial_device_monitor_error_string(serial_device_monitor_error err);

/* TODO: make this opaque <14-09-23, Max Schulte> */
typedef struct serial_device_monitor_context {
  struct udev *udev;
  struct udev_monitor *monitor;
} serial_device_monitor_context;

extern serial_device_monitor_error
serial_device_monitor_init(serial_device_monitor_context *context);

serial_device_monitor_error
serial_device_monitor_receive(serial_device_monitor_context *context,
                              struct udev_device **out);
extern void
serial_device_monitor_deinit(serial_device_monitor_context *context);

#ifdef __cplusplus
}
#endif

#endif

#ifdef SERIAL_DEVICE_MONITOR_IMPLEMENTATION

const char *
serial_device_monitor_error_string(serial_device_monitor_error err) {
  switch (err) {
  case serial_device_monitor_error_nil:
    return "<nil>";
  case serial_device_monitor_error_null_device:
    return "provided null device";
  case serial_device_monitor_error_udev_ctx_init:
    return "failed to create udev context";
  case serial_device_monitor_error_udev_mon_init:
    return "failed to create udev monitor";
  case serial_device_monitor_error_add_subsys_filter:
    return "failed to add subsystem device type "
           "filter to monitor";
  case serial_device_monitor_error_mon_enable_failed:
    return "failed to add enable monitor";
  /* case serial_device_monitor_error_device_recieve_failed: */
  /*   return "failed to add device filter to monitor"; */
  default:
    return "unkown error";
  }
}

serial_device_monitor_error
serial_device_monitor_init(serial_device_monitor_context *context) {
  context->udev = udev_new();
  if (!context->udev) {
    serial_device_monitor_deinit(context);
    return serial_device_monitor_error_udev_ctx_init;
  }
  context->monitor = udev_monitor_new_from_netlink(context->udev, "udev");
  if (!context->monitor) {
    serial_device_monitor_deinit(context);
    return serial_device_monitor_error_udev_mon_init;
  }
  if (udev_monitor_filter_add_match_subsystem_devtype(context->monitor, "tty",
                                                      NULL) != 0) {
    serial_device_monitor_deinit(context);
    return serial_device_monitor_error_add_subsys_filter;
  }
  if (udev_monitor_enable_receiving(context->monitor) != 0) {
    serial_device_monitor_deinit(context);
    return serial_device_monitor_error_mon_enable_failed;
  }
  return serial_device_monitor_error_nil;
}

serial_device_monitor_error
serial_device_monitor_receive(serial_device_monitor_context *context,
                              struct udev_device **out) {
  if (out == NULL)
    return serial_device_monitor_error_null_device;
  while (1) {
    struct udev_device *device = udev_monitor_receive_device(context->monitor);
    if (device == NULL)
      return serial_device_monitor_error_device_recieve_failed;
    const char *action = udev_device_get_action(device);
    if (strcmp(action, "add") != 0) {
      udev_device_unref(device);
      continue;
    }
    *out = device;
    return serial_device_monitor_error_nil;
  }
}

void serial_device_monitor_deinit(serial_device_monitor_context *context) {
  if (context->monitor != NULL)
    udev_monitor_unref(context->monitor);
  if (context->monitor != NULL)
    udev_monitor_unref(context->monitor);
  if (context->udev != NULL)
    udev_unref(context->udev);
}

#endif
