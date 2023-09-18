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

  char device_path[4096] = {0};
  err = serial_device_monitor_receive(&ctx, device_path);
  if (err) {
    serial_device_monitor_deinit(&ctx);
    fprintf(stderr, "%s\n", serial_device_monitor_error_string(err));
    return 1;
  }
  printf("%s\n", device_path);

  serial_device_monitor_deinit(&ctx);
  return 0;
}
