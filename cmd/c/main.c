
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

/* #include <fcntl.h> */
/* #include <libudev.h> */
/* #include <signal.h> */
/* #include <stdbool.h> */
/* #include <stdio.h> */
/* #include <stdlib.h> */
/* #include <string.h> */
/* #include <sys/stat.h> */
/* #include <unistd.h> */

/* static volatile sig_atomic_t keep_running = 1; */
/* static void sig_handler(int _) { */
/*   (void)_; */
/*   keep_running = 0; */
/* } */

/* void print_udev_list_entries(struct udev_list_entry *entries) { */
/*   while (entries != NULL) { */
/*     const char *key = udev_list_entry_get_name(entries); */
/*     const char *value = udev_list_entry_get_value(entries); */
/*     printf("\t%s: %s\n", key, value); */
/*     entries = udev_list_entry_get_next(entries); */
/*   } */
/* } */

/* bool is_tty(const char *devnode) { */
/*   struct stat st; */
/*   if (stat(devnode, &st) == -1) */
/*     return false; */
/*   if (!S_ISCHR(st.st_mode)) */
/*     return false; */
/*   int fd = open(devnode, O_RDONLY); */
/*   int rc = isatty(fd); */
/*   close(fd); */
/*   return rc; */
/* } */

/* int main() { */
/*   // Create a Udev context */
/*   struct udev *udev = udev_new(); */
/*   if (!udev) { */
/*     fprintf(stderr, "Failed to create udev context\n"); */
/*     return 1; */
/*   } */

/*   // Create a Udev monitor and filter for hotplug events */
/*   struct udev_monitor *monitor = udev_monitor_new_from_netlink(udev, "udev");
 */
/*   if (!monitor) { */
/*     fprintf(stderr, "Failed to create udev monitor\n"); */
/*     udev_unref(udev); */
/*     return 1; */
/*   } */

/*   if (udev_monitor_filter_add_match_subsystem_devtype(monitor, "tty", NULL)
 * != 0) { */
/*     fprintf(stderr, "Failed to add subsystem device type filter to
 * monitor\n"); */
/*     udev_unref(udev); */
/*     return 1; */
/* 	} */
/*   if (udev_monitor_enable_receiving(monitor) != 0) { */
/*     fprintf(stderr, "Failed to add enable monitor\n"); */
/*     udev_unref(udev); */
/*     return 1; */
/* 	} */

/*   printf("Listening for Udev hotplug events...\n"); */
/*   signal(SIGINT, sig_handler); */
/*   while (keep_running) { */
/*     struct udev_device *device = udev_monitor_receive_device(monitor); */
/*     if (device == NULL) */
/*       continue; */
/*     const char *action = udev_device_get_action(device); */
/*     if (strcmp(action, "add") != 0) { */
/*       udev_device_unref(device); */
/*       continue; */
/*     } */
/*     const char *devnode = udev_device_get_devnode(device); */
/*     if (!is_tty(devnode)) { */
/*       udev_device_unref(device); */
/*       continue; */
/*     } */
/*     printf("new terminal device: %s\n", devnode); */
/*     udev_device_unref(device); */
/*   } */

/*   printf("cleaining up\n"); */
/*   udev_monitor_unref(monitor); */
/*   udev_unref(udev); */
/*   return 0; */
/* } */
